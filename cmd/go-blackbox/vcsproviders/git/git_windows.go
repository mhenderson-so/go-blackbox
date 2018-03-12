// +build windows

package git

import (
	"os"
	"os/exec"
	"strings"

	"github.ds.stackexchange.com/mhenderson/go-blackbox/cmd/models"
)

// GetRepoBase returns the path to the base of the repository we're currently in
func (vcs *vcs) GetRepoBase() string {
	gitOut, _ := exec.Command("cmd", "/c", "git rev-parse --show-toplevel").Output()
	currentDir, _ := os.Getwd()

	pathCurrent := strings.TrimSpace(models.ConsistentSlashes(currentDir))
	pathGit := strings.TrimSpace(models.ConsistentSlashes(string(gitOut)))

	if string(gitOut) != "" && strings.Contains(pathCurrent, pathGit) {
		return pathGit
	}

	return ""
}
