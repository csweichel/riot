// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"

	"github.com/32leaves/riot/pkg/projectlib"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create (node|app) name",
	Short: "Creates a new resource in a cobra project",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			cmd.Help()
			return
		}

		basedir := getBaseDir(cmd)
		rt := args[0]
		name := args[1]
		var err error
		if rt == "node" {
			err = createNewNode(basedir, name)
		} else if rt == "app" {
			err = projectlib.CreateApplication(basedir, name)
		} else {
			err = fmt.Errorf("unknown resource type %s. Please use \"node\" or \"app\"", rt)
		}

		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func createNewNode(basedir string, name string) error {
	env, err := projectlib.LoadEnv(basedir)
	if err != nil {
		return err
	}

	env.AddNode(projectlib.Node{
		Name: name,
		Host: name,
	})
	return env.Save()
}
