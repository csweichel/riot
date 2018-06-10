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
	"io/ioutil"
	"net"
	"os"
	"path"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

// Environment is the core configuration of a riot project
type Environment interface {
	GetNodes() []Node
	GetApplications() ([]Application, error)
}

type environment struct {
	basedir string
	Nodes   []Node `yaml:"nodes"`
}

// Node represents a single device on which we can deploy an application to
type Node struct {
	Name   string   `yaml:"name"`
	Host   string   `yaml:"host"`
	Labels []string `yaml:"labels"`
}

// Nodes returns all nodes configured in an environment
func (env *environment) GetNodes() []Node {
	return env.Nodes
}

func (env *environment) GetApplications() ([]Application, error) {
	matches, err := filepath.Glob(path.Join(env.basedir, "applications", "*", "application.yaml"))
	if err != nil {
		return nil, err
	}

	result := make([]Application, len(matches))
	for idx, fn := range matches {
		appBasedir, _ := path.Split(fn)
		app, err := LoadApp(appBasedir)
		if err != nil {
			return nil, err
		}

		result[idx] = *app
	}
	return result, nil
}

// LoadEnv loads the environment description of a riot project from the basedir
func LoadEnv(basedir string) (Environment, error) {
	fn := path.Join(basedir, "environment.yaml")
	_, err := os.Stat(fn)
	if os.IsNotExist(err) {
		return nil, err
	}

	yamlFile, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	var result environment
	result.basedir = basedir
	err = yaml.Unmarshal(yamlFile, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// IsAvailable checks if a node is available for container deployment
func (node *Node) IsAvailable() bool {
	_, err := net.DialTimeout("tcp", node.Host+":2376", time.Duration(1)*time.Second)
	return err == nil
}
