// +build !windows

package git

// Set REPOBASE to the top of the repository
func (vcs *vcs) GetRepoBase(path string) string {
	//gitOut, err := exec.Command("which >/dev/null 2>/dev/null git && git rev-parse --show-toplevel >/dev/null 2>&1").CombinedOutput()
	return "" //fmt.Errorf("not implemented on Linux yet")
}
