package leaks

import "errors"

var (
	VirusDetectedErr = errors.New("Virus has been detected")
	FileNotFound     = errors.New("No such file")
	FileCheckErr     = errors.New("There was an error, while checking the file for viruses")
)
