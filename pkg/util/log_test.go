package util

import (
	"log"
	"testing"
)

// This test function demonstrates how to use logToFile
func Test_logToFile(t *testing.T) {
	file, err := logToFile()
	if err == nil {
		defer file.Close()

		log.SetOutput(file)
	}

	log.Println("This text will be printed in the testlogfile file")
}
