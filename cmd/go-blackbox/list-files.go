package blackbox

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

// ListFiles returns a []string of all the blackboxed files
// in the repository
func ListFiles() ([]string, error) {
	path := filepath.Join(blackboxData, "blackbox-files.txt")
	pathContents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// The TrimSpace here removes any trailing line endings, which means
	// we don't produce an extra empty element
	fileList := strings.Split(strings.TrimSpace(string(pathContents)), "\n")

	return fileList, nil
}
