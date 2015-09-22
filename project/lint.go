package project

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"

	gwl "github.com/adamveld12/gowatch/log"
	linter "github.com/golang/lint"
)

func lint(projectDirectory string) bool {
	gwl.LogDebug("linting code...")
	lint := &linter.Linter{}

	files := make(map[string]map[string][]byte)
	filepath.Walk(projectDirectory, func(p string, info os.FileInfo, err error) error {
		if filepath.Ext(p) == ".go" {
			fileWithPackage := strings.TrimPrefix(p, projectDirectory)
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
			linterConfidenceThresholdReached := false
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
					linterConfidenceThresholdReached = true
				} else {
					color.Yellow(lintLineOutput)
				}
			}
			lintErrors = linterConfidenceThresholdReached
		}
	}

	return !lintErrors
}
