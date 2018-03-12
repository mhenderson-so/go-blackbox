// +build !windows

package blackbox

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Set REPOBASE to the top of the repository
// Set VCS_TYPE to 'git', 'hg', 'svn' or 'unknown'
func setRepoBase() error {
	gitOut, err := exec.Command("which >/dev/null 2>/dev/null git && git rev-parse --show-toplevel >/dev/null 2>&1").CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error getting repository base path: %s", err)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Error getting current working directory: %s", err)
	}

	if strings.Contains(currentDir, string(gitOut)) {
		RepoBase = string(gitOut)
		VCSType = "git"
		return nil
	}

	return nil
}
