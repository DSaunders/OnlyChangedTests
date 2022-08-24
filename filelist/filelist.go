package filelist

import (
	"catalyst/config"
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"
)

type FileList struct {
	allFiles []string
	config   *config.Config
}

func BuildFor(rootPath string, config *config.Config) *FileList {
	list := FileList{
		allFiles: make([]string, 0),
		config:   config,
	}

	filepath.WalkDir(rootPath, func(s string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		for _, ignoredPath := range config.IgnorePaths {
			if d.IsDir() && strings.Contains(d.Name(), ignoredPath) {
				return filepath.SkipDir
			}
		}

		if !d.IsDir() {
			list.allFiles = append(list.allFiles, s)
		}

		return nil
	})

	return &list
}

func (filelist *FileList) FindTests() []string {

	regex := filelist.config.TestRegex

	tests := make([]string, 0)
	testRegex := regexp.MustCompile(regex)

	for _, file := range filelist.allFiles {
		if testRegex.MatchString(file) {
			tests = append(tests, file)
		}
	}

	return tests
}

func (filelist *FileList) Exists(filename string) bool {

	for _, file := range filelist.allFiles {
		if file == filename {
			return true
		}
	}

	return false
}
