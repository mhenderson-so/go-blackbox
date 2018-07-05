package git

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mhenderson-so/go-blackbox/cmd/models"
)

const (
	// RepoType is git
	RepoType = "git"
)

type vcs struct{}

func init() {
	git := &vcs{}
	models.RegisterVCS("git", git)
}

// GetRepoType returns git
func (vcs *vcs) GetRepoType() string {
	return RepoType
}

// Ignore adds the requested files/paths to the .gitignore file in the root of the Git repository
func (vcs *vcs) Ignore(repopath string, toignore []string) error {
	base := vcs.GetRepoBase(repopath)
	gitIgnore := filepath.Join(base, ".gitignore")

	return appendLines(gitIgnore, toignore)
}

func (vcs *vcs) Attributes(repopath string, attributes []string) error {
	base := vcs.GetRepoBase(repopath)
	gitAttributes := filepath.Join(base, ".gitattributes")

	return appendLines(gitAttributes, attributes)
}

// appendLines takes a file and checks to see if the current lines are already present.
// If the lines are not present, it appends them to the end of the file.
func appendLines(filepath string, things []string) error {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644) //Open our file
	if err != nil {
		return err
	}
	defer file.Close()

	//Read the contents of the file so we can check against it later.
	//Hope this isn't a big file or you might use all your RAM.
	currentContent, err := ioutil.ReadAll(file)
	currentContentString := string(currentContent)
	if err != nil {
		return err
	}

	//Go through each item in our list of things to append and check if they are in the file
	//already. If they are, don't do anything. If they're not, then append it to the end.
	//All items will have a newline attached to them, as that's how .git* files are built
	//(e.g. .gitignore, .gitattributes)
	for _, line := range things {
		exactLine := fmt.Sprintf("%s\n", line)
		if strings.Contains(currentContentString, exactLine) {
			continue
		}
		_, err := file.WriteString(exactLine)
		if err != nil {
			return err
		}
	}

	return nil
}
