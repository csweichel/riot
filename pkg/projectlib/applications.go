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
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// Application represnts a single app in a riot project
type Application struct {
	Name               string
	DeploymentSelector []string `yaml:"deploysTo"`
	Image              string   `yaml:"image"`
	BuildArgs          []string `yaml:"buildArgs"`
	RunArgs            []string `yaml:"runArgs"`
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
	var result []Node
	for _, c := range env.GetNodes() {
		nodeSelected := false
		for _, selector := range app.DeploymentSelector {
			if strings.HasPrefix(selector, "#") {
				// id selector
				if c.Name == strings.TrimPrefix(selector, "#") {
					result = append(result, c)
					nodeSelected = true
					break
				}
			} else if strings.HasPrefix(selector, ".") {
				selector := strings.TrimPrefix(selector, ".")
				labelFound := false
				for _, label := range c.Labels {
					if label == selector {
						labelFound = true
						break
					}
				}

				if labelFound {
					result = append(result, c)
					nodeSelected = true
					break
				}
			} else {
				return nil, fmt.Errorf("Invalid selector \"%s\". Must start with . or #", selector)
			}
		}

		if nodeSelected {
			break
		}
	}

	return result, nil
}
