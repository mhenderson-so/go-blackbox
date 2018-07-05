package blackbox

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mhenderson-so/go-blackbox/cmd/models"
)

// Initialize takes the path to a VCS repository and creates the
// required files and directories for working with BlackBox. This
// should be called on any new repository before attempting to
// use the rest of the BlackBox commands.
func Initialize(path string) error {
	vcs := models.GetActiveCVS(path)
	currentVCS := *vcs
	if currentVCS.GetRepoType() == "unknown" {
		return fmt.Errorf("not in a known VCS directory (%s)", path)
	}

	repoBase := currentVCS.GetRepoBase(path)

	blackboxPath := blackboxDataCandidates[0]
	blackboxDatapath := filepath.Join(repoBase, blackboxPath)

	//Stop Git from adding temporary or secret files
	toIgnore := []string{
		filepath.Join(blackboxPath, "pubring.gpg~"),
		filepath.Join(blackboxPath, "pubring.kbx~"),
		filepath.Join(blackboxPath, "secring.gpg"),
	}
	currentVCS.Ignore(repoBase, toIgnore)

	//Stop Windows from breaking the admin files
	toAttribute := []string{
		fmt.Sprintf("%s text eol=lf", blackboxAdminsFile),
		fmt.Sprintf("%s text eol=lf", blackboxFilesFile),
	}
	currentVCS.Attributes(repoBase, toAttribute)

	//Create the blackbox path
	err := os.MkdirAll(blackboxPath, 0666)
	if err != nil {
		return err
	}
	//Touch the admins file
	f, err := os.OpenFile(filepath.Join(blackboxDatapath, blackboxAdminsFile), os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	f.Close()

	//Touch the files file
	f, err = os.OpenFile(filepath.Join(blackboxDatapath, blackboxFilesFile), os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	f.Close()

	return nil
}
