package utils

import (
	"os"
)

// FileExist ...
func FileExist(f string) bool {

	if _, err := os.Stat(f); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}
