package unknown

import (
	"os"

	"github.ds.stackexchange.com/mhenderson/go-blackbox/cmd/models"
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

// GetRepoType returns git
func (vcs *vcs) GetRepoType() string {
	return RepoType
}
