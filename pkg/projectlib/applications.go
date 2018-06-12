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

package projectlib

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	yaml "gopkg.in/yaml.v2"
)

// Application represnts a single app in a riot project
type Application struct {
	Name               string
	DeploymentSelector []string `yaml:"deploysTo"`
	BuildCfg           AppBuild `yaml:"build"`
	Image              string   `yaml:"image"`
	RunCfg             AppRun   `yaml:"run"`
}

// AppBuild contains all settings related to an application build
type AppBuild struct {
	NodeSelector string             `yaml:"buildsOn"`
	Args         map[string]*string `yaml:"args"`
	DontPush     bool               `yaml:"dontPush"`
}

// AppRun configures an application during runtime
type AppRun struct {
	Volumes map[string]string `yaml:"volumes"`
	Ports   map[string]string `yaml:"ports"`
}

// LoadApp loads the application manifest from an application folder
func LoadApp(basedir string) (*Application, error) {
	fn := path.Join(basedir, "application.yaml")
	_, err := os.Stat(fn)
	if os.IsNotExist(err) {
		return nil, err
	}

	yamlFile, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	var result Application
	err = yaml.Unmarshal(yamlFile, &result)
	if err != nil {
		return nil, err
	}

	result.Name = path.Base(basedir)

	return &result, nil
}

// SelectDeploymentTargets selects all nodes in an environment to which an application ought to be deployed
func (app *Application) SelectDeploymentTargets(env Environment) ([]Node, error) {
	selectedNodes := make(map[string]Node)
	for _, selector := range app.DeploymentSelector {
		nodes, err := env.SelectNodes(selector)
		if err != nil {
			return nil, err
		} else if len(nodes) == 0 {
			return nil, fmt.Errorf("Selector \"%s\" did not match a node", selector)
		}

		for _, node := range nodes {
			selectedNodes[node.Name] = node
		}
	}

	result := make([]Node, 0)
	for _, c := range selectedNodes {
		result = append(result, c)
	}

	return result, nil
}

// IsDeployedOn checks if an application is currently deployed on a node
func (app *Application) IsDeployedOn(node Node) bool {
	// TODO: implement me
	return false
}

// GetBuildNode returns the node on which we should build the application image
func (app *Application) GetBuildNode(env Environment) (Node, error) {
	if len(app.BuildCfg.NodeSelector) > 0 {
		nodes, err := env.SelectNodes(app.BuildCfg.NodeSelector)
		if err != nil {
			return Node{}, err
		} else if len(nodes) == 0 {
			return Node{}, fmt.Errorf("Selector \"%s\" did not match a node", app.BuildCfg.NodeSelector)
		}

		return nodes[0], nil
	}

	nodes, err := app.SelectDeploymentTargets(env)
	if err != nil {
		return Node{}, err
	}

	return nodes[0], nil
}
