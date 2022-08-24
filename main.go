package main

import (
	"catalyst/config"
	"catalyst/dependencytree"
	"catalyst/filelist"
	"catalyst/git"
	"catalyst/jest"
	"catalyst/utils"
	"log"
	"strings"
)

func main() {

	config := config.LoadConfig("catalyst.config.json")

	changedFilesRelative := git.GetFilesChangedFromRevision(config.CompareToRevision)
	changedFilesAbsolute := utils.RelativeToAbsolutePaths(changedFilesRelative)
	log.Printf(
		"- Found %v changed files: \n  %v\n",
		len(changedFilesAbsolute),
		strings.Join(changedFilesAbsolute, "\n  "),
	)

	fileList := filelist.BuildFor(utils.GetExecutingPath(), config)

	log.Println("- Discovering test files")
	allTestFiles := fileList.FindTests()
	log.Printf("- Discovered %v test files\n", len(allTestFiles))

	log.Println("- Building dependency tree")
	depdencyTree := dependencytree.BuildForFiles(allTestFiles, fileList, config)

	log.Println("- Finding impacted tests")
	testsToRun := depdencyTree.GetTopLevelNodesForFiles(changedFilesAbsolute)
	log.Printf(
		"- There are %v impacted test files to run: \n  %v\n",
		len(testsToRun),
		strings.Join(testsToRun, "\n  "),
	)

	log.Println("- Generating Jest test selection code")
	jest.WriteFilterFile(testsToRun)

	log.Println("- Running Jest")
	jest.Run(config)
	log.Println("- Jest succeeded. All tests passed")
}
