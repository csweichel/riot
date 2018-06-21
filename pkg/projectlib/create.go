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
)

// CreateProject initiales a directory with the files and folders required for a riot project
func CreateProject(basedir string) error {
	fi, err := os.Stat(basedir)
	if os.IsNotExist(err) {
		err := os.Mkdir(basedir, os.ModePerm)
		if err != nil {
			return err
		}
	} else if !fi.IsDir() {
		return fmt.Errorf("Project path exists but is not a directory: %s", basedir)
	}

	err = ioutil.WriteFile(path.Join(basedir, "environment.yaml"), []byte(`registry:
  host: the-registry.local
nodes:
- name: myFirstNode
  host: my-first-node.local
  labels:
  - zerow
  - ble
  - concentrator`), 0644)
	if err != nil {
		return err
	}

	err = os.Mkdir(path.Join(basedir, "applications"), os.ModePerm)
	if err != nil {
		return err
	}

	err = CreateApplication(basedir, "with-build")
	if err != nil {
		return err
	}

	err = os.Mkdir(path.Join(basedir, "applications", "without-build"), os.ModePerm)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path.Join(basedir, "applications", "without-build", "application.yaml"), []byte(`deploysTo:
  - "#myFirstNode"
image: alpine:3.7`), 0644)
	if err != nil {
		return err
	}

	return nil
}

// CreateApplication creates a single application with a Dockerfile
func CreateApplication(basedir string, name string) error {
	err := os.Mkdir(path.Join(basedir, "applications", name), os.ModePerm)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path.Join(basedir, "applications", name, "Dockerfile"), []byte(`
FROM alpine
CMD ["echo", "hello"]
    `), 0644)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path.Join(basedir, "applications", name, "application.yaml"), []byte(`deploysTo:
  - ".ble"
build:
  buildsOn: ".ble"
  args:
    foo: bar
run:
  ports:
    8080: 8080
  volumes:
    /tmp/wbtmp: /tmp`), 0644)
	if err != nil {
		return err
	}

	return nil
}
