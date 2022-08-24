package filelist

import (
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"
)

type FileList struct {
	allFiles []string
}

func BuildFor(rootPath string) *FileList {
	list := FileList{
		allFiles: make([]string, 0),
	}

	filepath.WalkDir(rootPath, func(s string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// These should be configurable
		if d.IsDir() && strings.Contains(d.Name(), "node_modules") {
			return filepath.SkipDir
		}
		if d.IsDir() && strings.Contains(d.Name(), ".idea") {
			return filepath.SkipDir
		}
		if d.IsDir() && strings.Contains(d.Name(), "coverage") {
			return filepath.SkipDir
		}
		if d.IsDir() && strings.Contains(d.Name(), ".git") {
			return filepath.SkipDir
		}

		if !d.IsDir() {
			list.allFiles = append(list.allFiles, s)
		}
		return nil
	})

	return &list
}

func (filelist *FileList) FindTests() []string {

	// TODO: Get this from a config file
	regex := `(/_tests/.*|(\.|/)(test|spec))\.tsx?$`

	tests := make([]string, 0)
	testRegex := regexp.MustCompile(regex)

	for _, file := range filelist.allFiles {
		if testRegex.MatchString(file) {
			tests = append(tests, file)
		}
	}

	return tests
}
