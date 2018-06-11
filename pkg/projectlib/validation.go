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
	"strconv"
)

// Issue reports a single problem found in a project configuration
type Issue struct {
	Description string
	IsFatal     bool
}

func (issue Issue) String() string {
	level := "WARN"
	if issue.IsFatal {
		level = "ERROR"
	}
	return fmt.Sprintf("[%s] %s", level, issue.Description)
}

func (env *environment) Validate() ([]Issue, error) {
	result := make([]Issue, 0)

	issues, err := env.validateNodeNames()
	if err != nil {
		return nil, err
	} else {
		result = append(result, issues...)
	}

	issues, err = env.validateNodePorts()
	if err != nil {
		return nil, err
	} else {
		result = append(result, issues...)
	}

	return result, nil
}

func (env *environment) validateNodeNames() ([]Issue, error) {
	result := make([]Issue, 0)
	nodeNames := make(map[string]bool)
	for _, node := range env.Nodes {
		// TODO: find out how to check if key is in map
		if nodeNames[node.Name] {
			result = append(result, Issue{Description: "Node name is not unique: " + node.Name, IsFatal: true})
		}

		nodeNames[node.Name] = true
	}
	return result, nil
}

func (env *environment) validateNodePorts() ([]Issue, error) {
	nodeAppMap := make(map[string][]Application)
	apps, err := env.GetApplications()
	if err != nil {
		return nil, err
	}
	for _, app := range apps {
		nodes, err := app.SelectDeploymentTargets(env)
		if err != nil {
			return nil, err
		}

		for _, node := range nodes {
			nodeAppMap[node.Name] = append(nodeAppMap[node.Name], app)
		}
	}

	result := make([]Issue, 0)
	for nodeName, apps := range nodeAppMap {
		portsUsed := make(map[string]string)
		for _, app := range apps {
			for sourcePort, targetPort := range app.Ports {
				if nr, err := strconv.Atoi(sourcePort); err != nil || nr < 0 || nr > 65535 {
					result = append(result, Issue{
						Description: fmt.Sprintf("Source port %s on application %s is not a valid port number", targetPort, app.Name),
						IsFatal:     true,
					})
				}
				if nr, err := strconv.Atoi(targetPort); err != nil || nr < 0 || nr > 65535 {
					result = append(result, Issue{
						Description: fmt.Sprintf("Target port %s on application %s is not a valid port number", targetPort, app.Name),
						IsFatal:     true,
					})
				}

				if _, ok := portsUsed[targetPort]; ok {
					result = append(result, Issue{
						Description: fmt.Sprintf("Port %s on node %s is used by applications %s and %s", targetPort, nodeName, portsUsed[targetPort], app.Name),
						IsFatal:     true,
					})
				}
				portsUsed[targetPort] = app.Name
			}
		}
	}
	return result, nil
}
