package prjlayout

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
		err := os.Mkdir(basedir, os.ModeDir)
		if err != nil {
			return err
		}
	} else if !fi.IsDir() {
		return fmt.Errorf("Project path exists but is not a directory: %s", basedir)
	}

	err = ioutil.WriteFile(path.Join(basedir, "environment.yaml"), []byte(`
nodes:
  myFirstNode:
	mac: 00:80:41:ae:fd:7e
	labels:
	- zerow
	- ble
	- concentrator
	`), 0644)
	if err != nil {
		return err
	}

	err = os.Mkdir(path.Join(basedir, "applications"), os.ModeDir)
	if err != nil {
		return err
	}
	err = os.Mkdir(path.Join(basedir, "applications", "with-build"), os.ModeDir)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path.Join(basedir, "applications", "with-build", "Dockerfile"), []byte(`
FROM alpine
CMD ["echo", "hello"]
	`), 0644)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path.Join(basedir, "applications", "with-build", "application.yaml"), []byte(`
deploysTo:
  - ".ble"
buildArgs:
# - "-dockerBuildArgGoesHere"
runArgs:
# - "-dockerRunArgGoesHere"
	`), 0644)
	if err != nil {
		return err
	}

	err = os.Mkdir(path.Join(basedir, "applications", "without-build"), os.ModeDir)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path.Join(basedir, "applications", "without-build", "application.yaml"), []byte(`
deploysTo:
  - "#myFirstNode"
image: alpine:3.7
runArgs:
# - "-dockerRunArgGoesHere"
	`), 0644)
	if err != nil {
		return err
	}

	return nil
}
