package utils

import (
	"os"
	"path/filepath"
)

func GetExecutingPath() string {
	ex, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return ex
}

func RelativeToAbsolutePaths(rootPath string, relativePaths []string) []string {
	absolutePaths := make([]string, len(relativePaths))
	for i, path := range relativePaths {
		absolutePaths[i] = filepath.Join(rootPath, path)
	}
	return absolutePaths
}
