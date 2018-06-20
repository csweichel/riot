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
	"log"

	"github.com/32leaves/riot/pkg/projectlib"
	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build [app]",
	Short: "Builds applications of this project",
	Long: `Builds all (or the given) application docker images in this project and
pushes them to the main registry. After a successful build one can
deploy either the latest images (the last build) or a previous build.`,
	Run: func(cmd *cobra.Command, args []string) {
		basedir := getBaseDir(cmd)

		env, err := projectlib.LoadEnv(basedir)
		if err != nil {
			log.Fatal("Error while loading environment from ", basedir, "\n", err)
			return
		}

		var apps []projectlib.Application
        if len(args) > 0 {
            app, err := env.GetApplication(args[0])
            if err != nil {
                log.Fatal(err)
                return
            }
            apps = []projectlib.Application{app}
        } else {
            apps, err = env.GetApplications()
            if err != nil {
                log.Fatal("Error while loading application descriptions", err)
                return
            }
        }

		appToImage := make(map[string]string)
		for _, app := range apps {
			log.Printf("Building: %s\n", app.Name)
			iamgeName, err := app.Build(env)
			if err != nil {
				log.Fatalf("Error while building %s\n%s", app.Name, err)
				return
			}

			appToImage[app.Name] = iamgeName
		}

		lock := projectlib.RiotLock{Versions: appToImage}
		err = lock.Save(basedir)
		if err != nil {
			log.Fatal("Error while saving yarn lock: ", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
