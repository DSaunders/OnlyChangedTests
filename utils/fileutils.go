package utils

import (
	"os"
	"path/filepath"
)

func GetExecutingPath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)
}

func RelativeToAbsolutePaths(relativePaths []string) []string {
	currentFolder := GetExecutingPath()

	absolutePaths := make([]string, len(relativePaths))
	for i, path := range relativePaths {
		absolutePaths[i] = filepath.Join(currentFolder, path)
	}
	return absolutePaths
}
