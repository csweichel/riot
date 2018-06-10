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
	"github.com/fatih/color"
	"github.com/gosuri/uiprogress"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Displays the status of all applications and their deployment",
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

		nodes := env.GetNodes()
		uiprogress.Start()
		bar := uiprogress.AddBar(len(nodes)).AppendCompleted().PrependElapsed()
		bar.PrependFunc(func(b *uiprogress.Bar) string {
			return nodes[b.Current()-1].Name
		})
		hostAvailability := make([]bool, len(env.GetNodes()))
		for idx, node := range env.GetNodes() {
			bar.Incr()
			hostAvailability[idx] = node.IsAvailable()
		}

		apps, err := env.GetApplications()
		if err != nil {
			log.Fatal(err)
		}
		bar = uiprogress.AddBar(len(apps)).AppendCompleted().PrependElapsed()
		bar.PrependFunc(func(b *uiprogress.Bar) string {
			return apps[b.Current()-1].Name
		})
		applicationAvailability := make([]map[string]bool, len(apps))
		for idx, app := range apps {
			bar.Incr()

			applicationAvailability[idx] = make(map[string]bool)
			hosts, err := app.SelectDeploymentTargets(env)
			if err != nil {
				log.Fatal(err)
			}
			for _, node := range hosts {
				applicationAvailability[idx][node.Name] = app.IsDeployedOn(node)
			}
		}
		uiprogress.Stop()

		downColor := color.New(color.Bold, color.FgRed).SprintFunc()
		upColor := color.New(color.FgGreen).SprintFunc()
		for idx, node := range env.GetNodes() {
			var status string
			if hostAvailability[idx] {
				status = upColor("up")
			} else {
				status = downColor("down")
			}
			log.Printf("Host %s (node %s) is %s\n", node.Host, node.Name, status)
		}
		for idx, app := range apps {
			statement := app.Name + ":"
			for hn, status := range applicationAvailability[idx] {
				if status {
					statement += upColor(" +" + hn)
				} else {
					statement += downColor(" -" + hn)
				}
			}
			log.Println(statement)
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
