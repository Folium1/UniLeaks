package http

import (
	"fmt"
	"os"

	"leaks/pkg/models"
)

func createTempFile(leak models.Leak) (*os.File, error) {
	tempFile, err := os.CreateTemp("", leak.File.Id)
	if err != nil {
		return nil, err
	}
	if _, err := tempFile.Write(leak.File.Content); err != nil {
		return nil, err
	}
	return tempFile, nil
}

func deleteTempFile(tempFile *os.File) {
	if err := tempFile.Close(); err != nil {
		logger.Error(fmt.Sprint("Couldn't close file: ", err))
	}
	if err := os.Remove(tempFile.Name()); err != nil {
		logger.Error(fmt.Sprint("Couldn't remove file: ", err))
	}
}
