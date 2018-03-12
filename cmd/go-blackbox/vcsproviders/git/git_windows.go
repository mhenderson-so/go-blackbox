// +build windows

package git

import (
	"os/exec"
	"strings"

	"github.ds.stackexchange.com/mhenderson/go-blackbox/cmd/models"
)

// GetRepoBase returns the path to the base of the repository we're currently in
func (vcs *vcs) GetRepoBase(path string) string {
	gitOut, _ := exec.Command("cmd", "/c", "git", "-C", path, "rev-parse", "--show-toplevel").CombinedOutput()
	pathCurrent := strings.TrimSpace(models.ConsistentSlashes(path))
	pathGit := strings.TrimSpace(models.ConsistentSlashes(string(gitOut)))

	if string(gitOut) != "" && strings.Contains(pathCurrent, pathGit) {
		return pathGit
	}

	return ""
}
