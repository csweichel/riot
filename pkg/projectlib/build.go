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
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/docker/docker/api/types"
	"github.com/mholt/archiver"
	"github.com/rs/xid"
)

// Build builds the image of an application or returns the preconfigured one if there is no Dockerfile
func (app *Application) Build(env Environment) (string, error) {
	appBasedir := path.Join(env.GetBaseDir(), "applications", app.Name)
	log.Printf("Building application in %s", appBasedir)
	dockerfilePath := path.Join(appBasedir, "Dockerfile")
	if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
		return app.Image, nil
	}

	node, err := app.GetBuildNode(env)
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	client, err := node.GetDockerClient(ctx, env)
	if err != nil {
		return "", err
	}

	fileinfo, err := ioutil.ReadDir(appBasedir)
	if err != nil {
		return "", err
	}
	filelist := make([]string, 0)
	for _, fi := range fileinfo {
		if fi.Name() != "." && fi.Name() != ".." {
			filelist = append(filelist, path.Join(appBasedir, fi.Name()))
		}
	}
	tarfile, err := ioutil.TempFile("", "riot-build")
	if err != nil {
		return "", err
	}
	err = archiver.TarGz.Make(tarfile.Name(), filelist)
	if err != nil {
		return "", err
	}
	log.Printf("Built build-context at %s", tarfile.Name())

	imageVersion := xid.New().String()
	imageName := env.GetRegistry().Host + "/" + app.Name + ":" + imageVersion
	log.Printf("Building image %s\n", imageName)

	dockerBuildContext, err := os.Open(tarfile.Name())
	defer dockerBuildContext.Close()
	options := types.ImageBuildOptions{
		SuppressOutput: false,
		Remove:         true,
		ForceRemove:    true,
		PullParent:     true,
		BuildArgs:      app.BuildCfg.Args,
		Tags:           []string{imageName},
	}
	buildResponse, err := client.ImageBuild(ctx, dockerBuildContext, options)
	if err != nil {
		return "", err
	}
	defer buildResponse.Body.Close()
	io.Copy(os.Stdout, buildResponse.Body)

	if !app.BuildCfg.DontPush {
		authString, err := env.GetRegistry().GetAuthString()
		if err != nil {
			return "", err
		}
		pushOptions := types.ImagePushOptions{
			RegistryAuth: authString,
		}
		pushResponse, err := client.ImagePush(ctx, imageName, pushOptions)
		if err != nil {
			return "", err
		}
		defer pushResponse.Close()
		io.Copy(os.Stdout, pushResponse)
	}

	return imageName, nil
}
