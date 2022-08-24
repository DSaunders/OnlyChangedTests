package git

import (
	"log"
	"os/exec"
	"strings"
)

func GetFilesChangedFromRevision(revision string) []string {
	out, err := exec.Command("git", "diff", "--name-only", revision).Output()
	if err != nil {
		log.Fatal(err)
	}

	allFiles := strings.Split(strings.ReplaceAll(string(out), "\r\n", "\n"), "\n")

	withoutEmptyItems := make([]string, 0)

	for _, file := range allFiles {
		if file != "" {
			withoutEmptyItems = append(withoutEmptyItems, file)
		}
	}

	return withoutEmptyItems
}
