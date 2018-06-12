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
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

// RiotLock locks an application to a specific image version
type RiotLock struct {
	Versions map[string]string `yaml:"versions"`
}

// Save stores a riot lock in a project
func (lock RiotLock) Save(basedir string) error {
	data, err := yaml.Marshal(lock)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path.Join(basedir, "riot.lock"), data, os.ModePerm)
	return err
}

// LoadLock loads a riot.lock file for a project
func LoadLock(basedir string) (RiotLock, error) {
	fn := path.Join(basedir, "riot.lock")
	_, err := os.Stat(fn)
	if os.IsNotExist(err) {
		return RiotLock{}, err
	}

	yamlFile, err := ioutil.ReadFile(fn)
	if err != nil {
		return RiotLock{}, err
	}

	var result RiotLock
	err = yaml.Unmarshal(yamlFile, &result)
	if err != nil {
		return RiotLock{}, err
	}

	return result, nil
}
