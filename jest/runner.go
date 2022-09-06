package jest

import (
	"log"
	"onlychangedtests/config"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/TwiN/go-color"
)

func Run(config *config.Config) {
	log.Println(".. Running Jest")

	start := time.Now()

	cmd := exec.Command(config.JestPath, "--filter=./selected-tests.js")
	if err := cmd.Start(); err != nil {
		log.Fatalf(color.Ize(color.Red, ".. Couldn't start Jest: %v\n"), err)
	}

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if _, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				log.Printf(".. Jest took %v\n", time.Since(start))
				log.Println(color.Ize(color.Red, ".. Jest failed. There are likely test failures"))
				RemoveFilterFile()
				os.Exit(1)
			}
		} else {
			log.Fatalf(color.Ize(color.Red, ".. Couldn't start Jest: %v\n"), err)
			RemoveFilterFile()
		}
	}

	log.Printf(".. Jest took %v\n", time.Since(start))
	RemoveFilterFile()
}
