// +build !windows

package git

import (
	"os/exec"
	"strings"

	"github.com/mhenderson-so/go-blackbox/cmd/models"
)

// Set REPOBASE to the top of the repository
func (vcs *vcs) GetRepoBase(path string) string {
	//gitOut, err := exec.Command("which >/dev/null 2>/dev/null git && git rev-parse --show-toplevel >/dev/null 2>&1").CombinedOutput()
	gitOut, _ := exec.Command("git", "-C", path, "rev-parse", "--show-toplevel").CombinedOutput()
	pathCurrent := strings.TrimSpace(models.ConsistentSlashes(path))
	pathGit := strings.TrimSpace(models.ConsistentSlashes(string(gitOut)))

	if string(gitOut) != "" && strings.Contains(pathCurrent, pathGit) {
		return pathGit
	}

	return ""
}
