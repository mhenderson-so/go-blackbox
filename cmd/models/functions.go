package models

import (
	"os"
	"strings"
)

// ConsistentSlashes is used on Windows and Linux paths to ensure that all slashes across
// blackbox are consistent. This converts backslashes into forwardslashes (windows to unix style)
func ConsistentSlashes(path string) string {
	return strings.Replace(path, "\\", "/", -1)
}

// Exists - reference https://stackoverflow.com/a/10510718/69683
func Exists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, err
		}
		return true, err

	}

	return true, nil
}
