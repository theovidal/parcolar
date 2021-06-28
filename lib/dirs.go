package lib

import (
	"io/ioutil"
	"os"
)

var TempDir string
var WorkDir string

// OpenDirs loads current and temporary directories
func OpenDirs() {
	temp, err := ioutil.TempDir("", "parcolar-")
	if err != nil {
		Fatal("Error opening temporary directory: %s", err)
	}
	TempDir = temp
	err = os.Chdir(TempDir)
	if err != nil {
		Fatal("Error accessing to temporary directory: %s", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		Fatal("Error opening working directory: %s", err)
	}
	WorkDir = wd
}
