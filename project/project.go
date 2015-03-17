package project

import (
	"log"
	"os"
)

func GoFiles(cwd string) []string {
	fileInfo, err := os.Open(cwd)
	if err != nil {
		log.Fatal(err)
	}

	infos, err := fileInfo.Readdir(0)
	fileInfo.Close()
	if err != nil {
		log.Fatal(err)
	}

	// recursively read files
}

func Folders(cwd string) []string {
	// recursively read folders
}
