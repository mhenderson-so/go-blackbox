package unknown

import (
	"fmt"
	"os"

	"github.com/mhenderson-so/go-blackbox/cmd/models"
)

const (
	RepoType = "unknown"
)

type vcs struct{}

func init() {
	unknown := &vcs{}
	models.RegisterVCS("unknown", unknown)
}

// GetRepoBase returns current working directory
func (vcs *vcs) GetRepoBase(path string) string {
	wd, _ := os.Getwd()
	return wd
}

// GetRepoType returns unknown
func (vcs *vcs) GetRepoType() string {
	return RepoType
}

// Ignore returns unknown
func (vcs *vcs) Ignore(string, []string) error {
	return fmt.Errorf("unknown VCS provider")
}

// Atttribute returns unknown
func (vcs *vcs) Attributes(string, []string) error {
	return fmt.Errorf("unknown VCS provider")
}
