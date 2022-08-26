package logger

import (
	"fmt"
	"log"
	"time"
)

type logWriter struct {
	includeTimestamp bool
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	if writer.includeTimestamp {
		return fmt.Print(time.Now().UTC().Format("2006-01-02 15:04:05") + "  " + string(bytes))
	} else {
		return fmt.Print(string(bytes))
	}
}

func Init(includeTimestamp bool) {
	log.SetFlags(0)
	log.SetOutput(&logWriter{includeTimestamp: includeTimestamp})
}
