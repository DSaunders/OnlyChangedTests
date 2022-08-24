package jest

import (
	"catalyst/config"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func Run(config *config.Config) {
	cmd := exec.Command(config.JestPath, "--filter=./selected-tests.js")
	if err := cmd.Start(); err != nil {
		log.Fatalf("Couldn't start Jest: %v\n", err)
	}

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if _, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				log.Printf("- Jest failed. There are likely test failures\n")
				os.Exit(1)
			}
		} else {
			log.Fatalf("Couldn't start Jest: %v\n", err)
		}
	}
}
