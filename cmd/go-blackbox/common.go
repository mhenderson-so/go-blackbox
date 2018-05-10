package blackbox

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.ds.stackexchange.com/mhenderson/go-blackbox/cmd/go-blackbox/vcsproviders"
	"github.ds.stackexchange.com/mhenderson/go-blackbox/cmd/models"
)

var (
	// RepoBase is the base path of the repository that BlackBox is working in
	RepoBase, _ = os.Getwd()

	// VCSType one of 'git', 'hg', 'svn' or 'unknown', indicating what kind of
	// VCS repository we are working in
	VCSType = "unknown"

	blackboxDataCandidates = []string{
		filepath.Join("keyrings", "live"),
		".blackbox",
	}
	blackboxData = ""
	activeVCS    models.VCSProvider

	blackboxAdminsFile = "blackbox-admins.txt"
	blackboxFilesFile  = "blackbox-files.txt"
)

// InitBlackbox must be called before you interact with blackbox for the first time.
// It does things like detect the VCS, root directory and data paths for blackbox.
// If you are working on a specific file, you should provide the path to the file you're
// working on (do not include the file itself). You can run this through filepath.Dir if needed.
// If you are not working on a specific file, then this defaults to the current
// working directory.
func InitBlackbox(path string) error {
	if path == "" {
		path, _ = os.Getwd()
	}

	//We need to work off an absolute path in case we're being fed a filename that is
	//outside our current working directory
	if !filepath.IsAbs(path) {
		cwd, _ := os.Getwd()
		path = filepath.Clean(filepath.Join(cwd, path))
	}

	vcs := models.GetActiveCVS(path)
	activeVCS = *vcs
	VCSType = activeVCS.GetRepoType()
	RepoBase = activeVCS.GetRepoBase(path)
	err := setBlackboxData()
	if err != nil {
		return fmt.Errorf("Error with blackbox data: %s", err)
	}

	return nil
}

func setBlackboxData() error {
	for _, candidate := range blackboxDataCandidates {
		path := filepath.Join(RepoBase, candidate)
		pathOK, _ := models.Exists(path)
		if pathOK {
			blackboxData = path
			return nil
		}
	}

	return fmt.Errorf("Unable to detect blackbox data path. Have you run blackbox initialize?")
}

// IsOnCryptList returns whether or not a given filename is a valid file in the
// blackbox list of files
func IsOnCryptList(filename string) bool {
	files, _ := ListFiles()
	filename, err := FilenameRelativeToVCSRoot(filename)
	if err != nil {
		log.Fatal(err)
		return false
	}
	for _, cryptFile := range files {
		if cryptFile == filename {
			return true
		}
	}

	return false
}

// FilenameRelativeToVCSRoot takes a given filename and returns the filename path
// in context of the VCS Root
func FilenameRelativeToVCSRoot(filename string) (string, error) {
	fullPath := filename
	if !filepath.IsAbs(fullPath) {
		cwd, _ := os.Getwd()
		fullPath = filepath.Join(cwd, filename)
	}

	fullPath = models.ConsistentSlashes(filepath.Clean(fullPath)) //I know that filepath.Clean() makes slashes consistent, but they're not consistent throughout the library
	if !strings.Contains(strings.ToLower(fullPath), strings.ToLower(RepoBase)) {
		return "", fmt.Errorf("%s is not contained within the root directory of %s", filename, RepoBase)
	}

	relativePath := fullPath[len(RepoBase)+1:]
	if relativePath[:2] == "/./" {
		relativePath = relativePath[3:]
	}
	return relativePath, nil
}
