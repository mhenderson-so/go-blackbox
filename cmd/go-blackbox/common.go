package blackbox

import (
	"fmt"
	"os"
	"path"
	"strings"
)

var (
	// RepoBase is the base path of the repository that BlackBox is working in
	RepoBase, _ = os.Getwd()

	// VCSType one of 'git', 'hg', 'svn' or 'unknown', indicating what kind of
	// VCS repository we are working in
	VCSType = "unknown"

	blackboxDataCandidates = []string{"keyrings/live", ".blackbox"}
	blackboxData           = ""
)

func InitBlackbox() error {
	err := setRepoBase()
	if err != nil {
		return fmt.Errorf("Error with blackbox repository: %s", err)
	}
	err = setBlackboxData()
	if err != nil {
		return fmt.Errorf("Error with blackbox data: %s", err)
	}

	return nil
}

func setBlackboxData() error {
	for _, candidate := range blackboxDataCandidates {
		path := path.Join(RepoBase, candidate)
		pathOK, err := exists(path)
		if err != nil {
			return fmt.Errorf("Error detecting blackbox data path: %s", err)
		}
		if pathOK {
			blackboxData = path
			return nil
		}
	}

	return fmt.Errorf("Unable to detect blackbox data path. Have you run blackbox initialize?")
}

// https://stackoverflow.com/a/10510718/69683
func exists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, err
		}
		return true, err

	}

	return true, nil
}

// We want all paths between Windows and Linux to have forward slashes. Just makes
// things easier.
func consistentSlashes(path string) string {
	return strings.Replace(path, "\\", "/", -1)
}
