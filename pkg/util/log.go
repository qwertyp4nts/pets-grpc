package util

import (
	"fmt"
	"os"
)

func logToFile() (*os.File, error) {
	f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}

	return f, nil
}
