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
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

// Deploy installs an application on a node
func (app *Application) Deploy(node Node, env Environment, lock RiotLock) (*RiotLock, error) {
	ctx := context.Background()
	client, err := node.GetDockerClient(ctx, env)
	if err != nil {
		return nil, err
	}

	imageName, ok := lock.Versions[app.Name]
	if !ok {
		return nil, fmt.Errorf("application %s has no riot.lock entry. Please run riot build", app.Name)
	}

	out, err := client.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return nil, err
	}
	io.Copy(os.Stdout, out)

	containerID, ok := lock.GetDeployment(app.Name, node.Name)
	if ok {
		err := client.ContainerStop(ctx, containerID, nil)
		if err != nil {
			return nil, err
		}
	}

	resp, err := client.ContainerCreate(ctx, &container.Config{
		Image: imageName,
	}, nil, nil, "")
	if err != nil {
		return nil, err
	}

	err = client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return nil, err
	}

	lock.AddDeployment(app.Name, node.Name, resp.ID)

	return &lock, nil
}
