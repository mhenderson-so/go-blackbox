package models

import "log"

// VCSProvider is the interface that all VCS drivers must meet in order
// to function inside BlackBox.
type VCSProvider interface {
	GetRepoType() string
	GetRepoBase(string) string
	Ignore(string, []string) error
	Attributes(string, []string) error
}

// vcsTypes contains the name and the provider for the registered VCS providers
var vcsTypes = map[string]*VCSProvider{}

// RegisterVCS registers a VCS Provider for use with BlackBox
func RegisterVCS(name string, VCS VCSProvider) {
	if _, ok := vcsTypes[name]; ok {
		log.Fatalf("Cannot register the VCS %s multiple times", name)
	}
	vcsTypes[name] = &VCS
}

// GetVCS returns the VCS Provider given its string name.
func GetVCS(name string) *VCSProvider {
	return vcsTypes[name]
}

// GetActiveCVS returns the currently active VCS provider for the current working directory.
// If there is no valid provider, it returns the `unknown` provider.
func GetActiveCVS(path string) *VCSProvider {
	for name, vcs := range vcsTypes {
		if name == "unknown" {
			continue
		}
		thisVCS := *vcs
		base := thisVCS.GetRepoBase(path) != ""
		if base {
			return vcs

		}
	}

	return GetVCS("unknown")
}
