package main

import (
	"log"
	"path"
	"runchangedtests/config"
	"runchangedtests/dependencytree"
	"runchangedtests/filelist"
	"runchangedtests/git"
	"runchangedtests/jest"
	"runchangedtests/logger"
	"runchangedtests/utils"
	"strings"
	"time"

	"github.com/TwiN/go-color"
)

func main() {

	start := time.Now()

	configLocation := path.Join(utils.GetExecutingPath(), "runchangedtests.config.json")
	config := config.LoadConfig(configLocation)

	logger.Init(config.IncludeTimestampInOutput)

	changedFilesRelative := git.GetFilesChangedFromRevision(config.CompareToRevision)
	changedFilesAbsolute := utils.RelativeToAbsolutePaths(git.GetGitPath(), changedFilesRelative)
	log.Printf("Found "+color.Ize(color.Bold, "%v changed file(s)")+"\n", len(changedFilesAbsolute))
	if config.LogFileFileNames {
		log.Printf("These are:\n  %v\n", strings.Join(changedFilesAbsolute, "\n  "))
	}

	fileList := filelist.BuildFor(utils.GetExecutingPath(), config)

	log.Println("Looking for test files")
	allTestFiles := fileList.FindTests()
	log.Printf("Discovered %v total test file(s)\n", len(allTestFiles))

	log.Println("Building dependency tree")
	depdencyTree := dependencytree.BuildForFiles(allTestFiles, fileList, config)

	log.Println("Finding impacted tests")
	testsToRun := depdencyTree.GetTopLevelNodesForFiles(changedFilesAbsolute)

	log.Printf("There are "+color.Ize(color.Bold, "%v impacted test file(s) to run")+"\n", len(testsToRun))
	if config.LogFileFileNames {
		log.Printf("These are:\n  %v\n", strings.Join(testsToRun, "\n  "))
	}

	log.Printf("Took %v\n", time.Since(start))
	if len(testsToRun) > 0 {
		log.Printf("Running %v test file(s) with Jest:", len(testsToRun))
		jest.WriteFilterFile(testsToRun)
		jest.Run(config)
		log.Println(color.Ize(color.Green, "Jest succeeded. All tests passed"))
	}

	log.Println("All done!")
}
