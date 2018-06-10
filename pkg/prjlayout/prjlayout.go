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

	err = ioutil.WriteFile(path.Join(basedir, "hom.xml"), []byte(`
<?xml version="1.0" ?>
<hom>
	<room id="livingroom">
		<node id="ship01" class="pi-zero-w ble" />
		<sensor id="ba:aa:fe:2a:8e" name="Mijia Living Room" class="mijia" />
	</room>
</hom>
	`), 0644)
	if err != nil {
		return err
	}

	err = os.Mkdir(path.Join(basedir, "applications"), os.ModeDir)
	if err != nil {
		return err
	}
	err = os.Mkdir(path.Join(basedir, "applications", "helloworld"), os.ModeDir)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path.Join(basedir, "applications", "helloworld", "Dockerfile"), []byte(`
FROM alpine
CMD ["echo", "hello"]
	`), 0644)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path.Join(basedir, "applications", "helloworld", "application.yaml"), []byte(`
deploysTo:
  - ".ble"
  - "#livingroom"
buildArgs:
# - "-dockerBuildArgGoesHere"
runArgs:
# - "-dockerRunArgGoesHere"
	`), 0644)
	if err != nil {
		return err
	}

	return nil
}
