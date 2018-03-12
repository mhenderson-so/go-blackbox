package blackbox

import (
	"fmt"
	"os"
	"path"

	_ "github.ds.stackexchange.com/mhenderson/go-blackbox/cmd/go-blackbox/vcsproviders"
	"github.ds.stackexchange.com/mhenderson/go-blackbox/cmd/models"
)

var (
	// RepoBase is the base path of the repository that BlackBox is working in
	RepoBase, _ = os.Getwd()

	// VCSType one of 'git', 'hg', 'svn' or 'unknown', indicating what kind of
	// VCS repository we are working in
	VCSType = "unknown"

	blackboxDataCandidates = []string{"keyrings/live", ".blackbox"}
	blackboxData           = ""
	activeVCS              models.VCSProvider
)

// InitBlackbox must be called before you interact with blackbox for the first time.
// It does things like detect the VCS, root directory and data paths for blackbox.
func InitBlackbox() error {
	vcs := models.GetActiveCVS()
	activeVCS = *vcs
	RepoBase = activeVCS.GetRepoBase()
	err := setBlackboxData()
	if err != nil {
		return fmt.Errorf("Error with blackbox data: %s", err)
	}

	return nil
}

func setBlackboxData() error {
	for _, candidate := range blackboxDataCandidates {
		path := path.Join(RepoBase, candidate)
		pathOK, _ := models.Exists(path)
		if pathOK {
			blackboxData = path
			return nil
		}
	}

	return fmt.Errorf("Unable to detect blackbox data path. Have you run blackbox initialize?")
}
