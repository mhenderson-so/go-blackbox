package blackbox

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

// ListAdmins returns a []string of all the blackbox administrators for the repository
func ListAdmins() ([]string, error) {
	path := filepath.Join(blackboxData, blackboxAdminsFile)
	pathContents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// The TrimSpace here removes any trailing line endings, which means
	// we don't produce an extra empty element
	fileList := strings.Split(strings.TrimSpace(string(pathContents)), "\n")

	return fileList, nil
}
