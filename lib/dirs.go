package lib

import (
	"io/ioutil"
	"log"
	"os"
)

var TempDir string
var WorkDir string

func OpenDirs() {
	temp, err := ioutil.TempDir("", "bacbot")
	if err != nil {
		log.Fatalln("Error opening temporary directory: " + err.Error())
	}
	TempDir = temp
	err = os.Chdir(TempDir)
	if err != nil {
		log.Panicln("Error accessing to temporary directory: " + err.Error())
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln("Error opening working directory: " + err.Error())
	}
	WorkDir = wd
}
