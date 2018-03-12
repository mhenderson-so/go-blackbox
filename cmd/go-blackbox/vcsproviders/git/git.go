package git

import "github.ds.stackexchange.com/mhenderson/go-blackbox/cmd/models"

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
