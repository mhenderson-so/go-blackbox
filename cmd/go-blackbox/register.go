package blackbox

import (
	"fmt"
	"path/filepath"

	"github.ds.stackexchange.com/mhenderson/go-blackbox/cmd/models"
)

// RegisterNewFile takes previously unencrypted file and enrolls it into
// the blackbox system.  Each file will be kept in the repo as an encrypted file.
func RegisterNewFile(filename string) error {
	gpgFilename := fmt.Sprintf("%s.gpg", filename)
	//Check if this file is already registered. If it is, don't register
	//it again.
	currentFiles, err := ListFiles()
	if err != nil {
		return err
	}

	//Check if this file is already in the list of blackbox files. We resolve the
	//absolute path to the file as this takes care of any issues with how the user
	//may have typed the command. For example adding ".\testfile.txt" would fail a
	//straight comparison, as the filename in currentFiles would be "testfile.txt".
	//By running these through filepath.Join(RepoBase,x) first we get a valid
	//comparison.
	absFilename := filepath.Join(RepoBase, filename)
	for _, current := range currentFiles {
		current = filepath.Join(RepoBase, current)
		if current == absFilename {
			return fmt.Errorf("Only register unencrypted files")
		}
	}

	//Check that this file actually exists on disk
	exists, err := models.Exists(filename)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("Only register files that exist")
	}

	//Check that the gpg version of this file does not already exist
	exists, err = models.Exists(gpgFilename)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("%s already exists", gpgFilename)
	}

	return nil
}
