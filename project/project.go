package project

import (
	"log"
	"os"
	"strings"
)

func GoFiles(cwd string) []string {
	dirHandle, err := os.Open(cwd)
	if err != nil {
		log.Fatal(err)
	}

	fileInfos, err := dirHandle.Readdir(0)
	dirHandle.Close()
	if err != nil {
		log.Fatal(err)
	}

	results := []string{}

	// recursively read files
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			for _, childFile := range GoFiles(fileInfo.Name()) {
				results = append(results, childFile)
			}
		} else if strings.HasSuffix(fileInfo.Name(), ".go") {
			results = append(results, fileInfo.Name())
		}
	}

	return results
}
