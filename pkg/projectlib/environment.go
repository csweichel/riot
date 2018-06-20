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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"gopkg.in/yaml.v2"
    "github.com/sahilm/fuzzy"
)

// Environment is the core configuration of a riot project
type Environment interface {
	GetRegistry() RegistryCfg
	GetNodes() []Node
	GetApplications() ([]Application, error)
    GetApplication(name string) (Application, error)
	Validate() ([]Issue, error)
	GetBaseDir() string
	SelectNodes(selector string) ([]Node, error)
}

type environment struct {
	basedir  string
	Registry RegistryCfg `yaml:"registry"`
	Nodes    []Node      `yaml:"nodes"`
}

// RegistryCfg configures access to a docker registry
type RegistryCfg struct {
	Host     string `yaml:"host"`
	Username string `yaml:"user"`
	Password string `yaml:"password"`
}

// Node represents a single device on which we can deploy an application to
type Node struct {
	Name   string   `yaml:"name"`
	Host   string   `yaml:"host"`
	Labels []string `yaml:"labels"`
}

// GetAuthString computes the base64 authorization string needed for docker registry requests
func (reg RegistryCfg) GetAuthString() (string, error) {
	auth := types.AuthConfig{
		Username: reg.Username,
		Password: reg.Password,
	}
	authBytes, err := json.Marshal(auth)
	if err != nil {
		return "", err
	}

	authBase64 := base64.URLEncoding.EncodeToString(authBytes)
	return authBase64, nil
}

func (env *environment) GetRegistry() RegistryCfg {
	return env.Registry
}

// Nodes returns all nodes configured in an environment
func (env *environment) GetNodes() []Node {
	return env.Nodes
}

func (env *environment) GetBaseDir() string {
	return env.basedir
}

func (env *environment) GetApplications() ([]Application, error) {
	matches, err := filepath.Glob(filepath.Join(env.basedir, "applications", "*", "application.yaml"))
	if err != nil {
		return nil, err
	}

	result := make([]Application, len(matches))
	for idx, fn := range matches {
		appBasedir, _ := filepath.Split(fn)
		app, err := LoadApp(appBasedir)
		if err != nil {
			return nil, err
		}

		result[idx] = *app
	}
	return result, nil
}

func (env *environment) GetApplication(name string) (Application, error) {
    applications, err := env.GetApplications()
    if err != nil {
        return Application{}, err
    }

    names := make([]string, 0)
    for _, app := range applications {
        if app.Name == name {
            return app, nil
        } else {
            names = append(names, app.Name)
        }
    }

    matches := fuzzy.Find(name, names)
    if len(matches) > 0 {
        return Application{}, fmt.Errorf("Application %s not found. Did you mean %s?", name, matches[0].Str)
    } else {
        return Application{}, fmt.Errorf("Application %s not found", name)
    }
}

// LoadEnv loads the environment description of a riot project from the basedir
func LoadEnv(basedir string) (Environment, error) {
	fn := filepath.Join(basedir, "environment.yaml")
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

func (env *environment) SelectNodes(selector string) ([]Node, error) {
	result := make([]Node, 0)
	for _, c := range env.Nodes {
		if strings.HasPrefix(selector, "#") {
			// id selector
			if c.Name == strings.TrimPrefix(selector, "#") {
				result = append(result, c)
				return result, nil
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
			}
		} else {
			return nil, fmt.Errorf("Invalid selector \"%s\". Must start with . or #", selector)
		}
	}
	return result, nil
}
