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

// vetCmd represents the init command
var vetCmd = &cobra.Command{
	Use:   "vet",
	Short: "Validates a riot project",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		basedir, err := rootCmd.PersistentFlags().GetString("project")
		if err != nil {
			log.Fatal(err)
			basedir = "."
		}

		env, err := projectlib.LoadEnv(basedir)
		if err != nil {
			log.Fatal("Error while loading environment from ", basedir, "\n", err)
			return
		}

		issues, err := env.Validate()
		if err != nil {
			log.Fatal("Error while vetting: ", err)
		} else {
			fatal := false
			for _, issue := range issues {
				fmt.Println(issue)
				fatal = fatal || issue.IsFatal
			}

			if fatal {
				log.Fatalf("Found fatal errors")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(vetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// vetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// vetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
