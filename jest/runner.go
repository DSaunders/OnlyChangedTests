package jest

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"runchangedtests/config"
	"syscall"
	"time"

	"github.com/TwiN/go-color"
)

func Run(config *config.Config) {
	log.Println(".. Running Jest")

	start := time.Now()

	cmd := exec.Command(config.JestPath, "--filter=./selected-tests.js", "--passWithNoTests", "--runInBand")

	// Create a buffer to stream the command output to stdout
	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)
	cmd.Stdout = mw
	cmd.Stderr = mw

	log.Println(color.Ize(color.Gray, ".. Jest output begins below"))
	log.Println(color.Ize(color.Gray, "-----"))

	// Run the command and stream the output
	err := cmd.Run()
	log.Println(stdBuffer.String())

	log.Println(color.Ize(color.Gray, "-----"))
	log.Println(color.Ize(color.Gray, ".. End of Jest output"))

	// Check the result of the command
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if _, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				log.Printf(".. Jest took %v\n", time.Since(start))
				log.Println(color.Ize(color.Red, ".. Jest failed. There are likely test failures"))
				log.Println(exiterr.Error())
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
