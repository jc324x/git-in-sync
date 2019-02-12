// Package stat tracks stats as the run progress.
package stat

import (
	"github.com/jychri/git-in-sync/pkg/brf"
)

// Stat tracks stats for the current run.
type Stat struct {
	CreatedWorkspaces      []string
	VerifiedWorkspaces     []string
	InaccessibleWorkspaces []string
	TotalWorkspaces        []string
	PendingClones          []string
	ClonedRepos            []string
	TotalRepos             []string
}

// Init returns a new *Stat.
func Init() *Stat {
	st := new(Stat)
	return st
}

// Reduce reduces slices in *Stat to their unique elements.
func (st *Stat) Reduce() {
	st.CreatedWorkspaces = brf.Reduce(st.CreatedWorkspaces)
	st.VerifiedWorkspaces = brf.Reduce(st.VerifiedWorkspaces)
	st.InaccessibleWorkspaces = brf.Reduce(st.InaccessibleWorkspaces)
	st.TotalWorkspaces = brf.Reduce(st.TotalWorkspaces)
}
