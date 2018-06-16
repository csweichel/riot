// Copyright © 2018 Christian Weichel <christian@csweichel.de>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"log"

	"github.com/32leaves/riot/pkg/projectlib"
	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys all applications of this project",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		basedir := getBaseDir(cmd)

		env, err := projectlib.LoadEnv(basedir)
		if err != nil {
			log.Fatal("Error while loading environment from ", basedir, "\n", err)
			return
		}

		apps, err := env.GetApplications()
		if err != nil {
			log.Fatal("Error while loading application descriptions", err)
			return
		}

		lock, err := projectlib.LoadLock(env.GetBaseDir())
		if err != nil {
			log.Fatal(err, ". Please run riot build.")
			return
		}

		errors := make([]error, 0)
		for _, app := range apps {
			hosts, err := app.SelectDeploymentTargets(env)
			if err != nil {
				errors = append(errors, err)
				continue
			}

			for _, host := range hosts {
				plock, err := app.Deploy(host, env, lock)

				if err != nil {
					errors = append(errors, err)
				} else {
					lock = *plock
					lock.Save(basedir)
				}
			}
		}

		if len(errors) > 0 {
			errorMessages := ""
			for _, err := range errors {
				errorMessages += fmt.Sprintf("%s\n", err)
			}
			log.Fatalf("Error while deploying project: %s", errorMessages)
		}
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
