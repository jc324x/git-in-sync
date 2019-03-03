// Package stat tracks stats as the run progress.
package stat

import (
	"github.com/jychri/goku/brf"
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

// Clear clears out slice based stats
func (st *Stat) Clear() {
	st.Workspaces = nil
	st.CreatedWorkspaces = nil
	st.VerifiedWorkspaces = nil
	st.InaccessibleWorkspaces = nil
	st.PendingClones = nil
	st.ClonedRepos = nil
	st.PendingRepos = nil
	st.ScheduledRepos = nil
	st.SkippedRepos = nil
	st.CompleteRepos = nil
	st.ScheduledPull = nil
	st.ScheduledPush = nil
}

// CheckComplete ...
func (st *Stat) CheckComplete() bool {
	switch {
	case len(st.Repos) == len(st.CompleteRepos):
		return true
	default:
		return false
	}
}

// CheckPending ...
func (st *Stat) CheckPending() bool {
	switch {
	case len(st.PendingRepos) >= 1:
		return true
	default:
		return false
	}
}

// CheckSkipped ...
func (st *Stat) CheckSkipped() bool {
	switch {
	case len(st.SkippedRepos) >= 1:
		return true
	default:
		return false
	}
}

// CheckScheduled ...
func (st *Stat) CheckScheduled() bool {
	switch {
	case len(st.SkippedRepos) >= 1:
		return true
	default:
		return false
	}
}
