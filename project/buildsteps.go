package project

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	gwl "github.com/adamveld12/gowatch/log"
	"github.com/fatih/color"

	linter "github.com/golang/lint"
)

func build(projectDirectory string, outputName string) bool {
	gwl.Debugln("\tbuilding", outputName)
	gwl.Debugln("\t@ dir", projectDirectory)
	cmd := exec.Command("go", "build", "-o", outputName)
	cmd.Dir = projectDirectory
	cmd.Env = os.Environ()

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return false
	}

	return true
}

func test(projectDirectory string) bool {
	gwl.LogDebug("testing code...")

	cmd := exec.Command("go", "test")
	cmd.Dir = projectDirectory
	cmd.Env = os.Environ()

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		gwl.LogDebug("test failures: ", err)
		return false
	}

	return true
}

func lint(projectDirectory string) bool {
	gwl.LogDebug("linting code...")
	lint := &linter.Linter{}

	files := walkFilesForLinting(projectDirectory)

	lintErrors := false
	for k, v := range files {
		gwl.LogDebug("linting package %s", k)

		problems, err := lint.LintFiles(v)

		if err != nil {
			color.Red("[ERROR]", err)
			lintErrors = true
		} else if len(problems) > 0 {

			gwl.LogDebug("lint issues found")
			color.Yellow("%d lint issue(s) found in %s\n\n", len(problems), k)
			for i, p := range problems {
				position := p.Position
				fileWithPackage := strings.Trim(strings.TrimPrefix(position.Filename, projectDirectory), "/")
				lintInfo := strings.Split(p.String(), "\n")

				gwl.LogDebug("%d out of 3", len(lintInfo))

				readableLintError := ""

				if len(lintInfo) >= 3 {
					readableLintError = fmt.Sprintf("- %s", lintInfo[2])
				}

				lintLineOutput := fmt.Sprintf("\t%d. %s line %d %s\n\t%s\n\n",
					i+1, fileWithPackage, position.Line, readableLintError, lintInfo[0])

				if p.Confidence > 0.5 {
					color.Red(lintLineOutput)
					lintErrors = true
				} else {
					color.Yellow(lintLineOutput)
				}
			}
		}
	}

	return !lintErrors
}

func walkFilesForLinting(packagePath string) map[string]map[string][]byte {
	files := make(map[string]map[string][]byte)
	filepath.Walk(packagePath, func(p string, info os.FileInfo, err error) error {
		if filepath.Ext(p) == ".go" {
			fileWithPackage := strings.TrimPrefix(p, packagePath)
			packageName := strings.Trim(strings.TrimSuffix(fileWithPackage, info.Name()), "/")

			if packageName == "" {
				packageName = "main"
			}

			files[packageName] = make(map[string][]byte)

			f, err := os.Open(p)
			if err != nil {
				return err
			}

			if files[packageName][p], err = ioutil.ReadAll(f); err != nil {
				return err
			}
		}
		return nil
	})

	return files
}
