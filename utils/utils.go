package utils

import (
	"os"
)

type Flags struct {
	Thread     bool
	ConfigPath string
}

func CheckFile(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}
