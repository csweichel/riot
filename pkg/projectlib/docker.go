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
	"log"
	"os"
	"path"

	"github.com/docker/docker/client"
)

// GetDockerClient builds a Docker client for a particular node
func (node *Node) GetDockerClient(ctx context.Context, env Environment) (*client.Client, error) {
	certs := path.Join(env.GetBaseDir(), ".riot-certs", node.Name)
	var cli *client.Client
	var err error
	if _, err := os.Stat(certs); os.IsNotExist(err) {
		log.Printf("Unable to find Docker certificates for node %s. Please run riot collect", node.Name)
		cli, err = client.NewClientWithOpts(
			client.WithVersion("1.37"),
			client.WithHost(node.getDockerHost()),
		)
	} else {
		cli, err = client.NewClientWithOpts(
			client.WithTLSClientConfig(
				path.Join(certs, "ca.pem"),
				path.Join(certs, "cert.pem"),
				path.Join(certs, "key.pem")),
			client.WithVersion("1.37"),
			client.WithHost(node.getDockerHost()),
		)
	}
	if err != nil {
		return nil, err
	}

	return cli, nil
}

func (node *Node) getDockerHost() string {
	return fmt.Sprintf("tcp://%s:2376", node.Host)
}
