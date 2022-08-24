package main

import (
	"dependencytracker/dependencytree"
	"dependencytracker/filelist"
	"dependencytracker/git"
	"dependencytracker/utils"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func main() {

	changedFilesRelative := git.GetFilesChangedFromRevision("origin/master")
	changedFilesAbsolute := utils.RelativeToAbsolutePaths(changedFilesRelative)
	log.Printf("\n\n%v changed files: \n - %v\n\n\n", len(changedFilesAbsolute), strings.Join(changedFilesAbsolute, "\n - "))

	currentFolder := utils.GetExecutingPath()
	fileList := filelist.BuildFor(currentFolder)

	allTests := fileList.FindTests()

	depdencyTree := dependencytree.BuildForFiles(allTests)

	testsToRun := depdencyTree.GetTopLevelNodesForFiles(changedFilesAbsolute)

	log.Printf("\n\n%v impacted test files to run: \n - %v\n\n\n", len(testsToRun), strings.Join(testsToRun, "\n - "))

	// Write the jest file
	err := ioutil.WriteFile("selected-tests.js", []byte(getTestFileContents(testsToRun)), 0644)
	if err != nil {
		panic(err)
	}

	log.Printf("Running jest\n")

	cmd := exec.Command(`.\node_modules\.bin\jest`, "--filter=./selected-tests.js")
	if err := cmd.Start(); err != nil {
		log.Fatalf("cmd.Start: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0

			// This works on both Unix and Windows. Although package
			// syscall is generally platform dependent, WaitStatus is
			// defined for both Unix and Windows and in both cases has
			// an ExitStatus() method with the same signature.
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				log.Printf("Exit Status: %d", status.ExitStatus())
				os.Exit(1)
			}
		} else {
			log.Fatalf("cmd.Wait: %v", err)
		}
	}
}

func getTestFileContents(files []string) string {
	code :=
		`
	const toRun = [
{{TESTS}}
	]

	module.exports = testPaths => {		
		return {
			filtered: testPaths.filter(a => {
				return toRun.includes(a.toLowerCase())
				})
				.map((testPath) => ({ test: testPath }))
			};
	};
`

	filesForJs := make([]string, len(files))
	for i, file := range files {
		filesForJs[i] = "        \"" + strings.Replace(file, "\\", "\\\\", -1) + "\","
	}

	code = strings.Replace(code, "{{TESTS}}", strings.Join(filesForJs, "\n"), 1)
	return code
}
