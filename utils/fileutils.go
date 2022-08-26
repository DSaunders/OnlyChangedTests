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

func RelativeToAbsolutePaths(relativePaths []string) []string {
	currentFolder := GetExecutingPath()

	absolutePaths := make([]string, len(relativePaths))
	for i, path := range relativePaths {
		absolutePaths[i] = filepath.Join(currentFolder, path)
	}
	return absolutePaths
}
