// +build windows

package blackbox

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

//go:generate goversioninfo -icon=build/archive_black_box_wLD_icon.ico

// Set REPOBASE to the top of the repository
// Set VCS_TYPE to 'git', 'hg', 'svn' or 'unknown'
func setRepoBase() error {
	gitOut, err := exec.Command("cmd", "/c", "git rev-parse --show-toplevel").Output()
	if err != nil {
		return fmt.Errorf("Error getting repository base path: %s", err)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Error getting current working directory: %s", err)
	}

	pathCurrent := strings.TrimSpace(consistentSlashes(currentDir))
	pathGit := strings.TrimSpace(consistentSlashes(string(gitOut)))

	if string(gitOut) != "" && strings.Contains(pathCurrent, pathGit) {
		RepoBase = pathGit
		VCSType = "git"
		return nil
	}

	return nil
}
