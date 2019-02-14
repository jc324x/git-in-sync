// Package stat tracks stats as the run progress.
package stat

import (
	"github.com/jychri/git-in-sync/pkg/brf"
)

// Stat tracks stats for the current run.
type Stat struct {
	Workspaces             []string
	CreatedWorkspaces      []string
	VerifiedWorkspaces     []string
	InaccessibleWorkspaces []string
	PendingClones          []string
	Repos                  []string
	ClonedRepos            []string
	PendingRepos           []string
	ScheduledRepos         []string
	SkippedRepos           []string
	CompleteRepos          []string
	ScheduledPull          []string
	ScheduledPush          []string
	// Complete               bool
}

// Init returns a new *Stat.
func Init() *Stat {
	st := new(Stat)
	return st
}

// Reduce reduces slices in *Stat to their unique elements.
func (st *Stat) Reduce() {
	st.Workspaces = brf.Reduce(st.Workspaces)
	st.CreatedWorkspaces = brf.Reduce(st.CreatedWorkspaces)
	st.VerifiedWorkspaces = brf.Reduce(st.VerifiedWorkspaces)
	st.InaccessibleWorkspaces = brf.Reduce(st.InaccessibleWorkspaces)
}

// IsComplete ...
func (st *Stat) IsComplete() bool {

	if len(st.Repos) == len(st.CompleteRepos) {
		return true

	}
	return false
}

// OnlyPending ...
func (st *Stat) OnlyPending() bool {
	return true
}

// OnlySkipped ...
func (st *Stat) OnlySkipped() bool {
	return true
}

// OnlyScheduled ...
func (st *Stat) OnlyScheduled() bool {
	return true
}
