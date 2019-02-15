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

// Continue ...
func (st *Stat) Continue() bool {
	switch {
	case st.AllComplete():
		return false
	case st.OnlySkipped():
		return false
	case len(st.ScheduledRepos) >= 1:
		return true
	default:
		return true
	}
}

// AllComplete ...
func (st *Stat) AllComplete() bool {
	if len(st.Repos) == len(st.CompleteRepos) {
		return true
	}

	return false
}

// OnlyPending ...
func (st *Stat) OnlyPending() bool {

	if len(st.SkippedRepos) >= 1 {
		return false
	}

	if len(st.ScheduledRepos) >= 1 {
		return false
	}

	return true
}

// OnlySkipped ...
func (st *Stat) OnlySkipped() bool {

	if len(st.PendingRepos) >= 1 {
		return false
	}

	if len(st.ScheduledRepos) >= 1 {
		return false
	}

	return true
}

// OnlyScheduled ...
func (st *Stat) OnlyScheduled() bool {

	if len(st.PendingRepos) >= 1 {
		return false
	}

	if len(st.SkippedRepos) >= 1 {
		return false
	}

	return true
}

// if scr := len(st.ScheduledPull); scr != 0 {
// 	etr := emoji.Get("Arrival")
// 	srs := brf.Summary(st.ScheduledPull, 12)
// 	brf.Printv(f, "%v [%v](%v) pull scheduled", etr, scr, srs)
// }

// if scr := len(st.ScheduledPush); scr != 0 {
// 	etr := emoji.Get("Departure")
// 	srs := brf.Summary(st.ScheduledPush, 12)
// 	brf.Printv(f, "%v [%v](%v) push scheduled", etr, scr, srs)
// }

// if pr := len(st.PendingRepos); pr != 0 {
// 	etr := emoji.Get("Traffic")
// 	srs := brf.Summary(st.PendingRepos, 12)
// 	brf.Printv(f, "%v [%v](%v) pending", etr, pr, srs)
// }

// if skr := len(st.SkippedRepos); skr != 0 {
// 	etr := emoji.Get("Stop")
// 	srs := brf.Summary(st.SkippedRepos, 12)
// 	brf.Printv(f, "%v [%v](%v) skipped", etr, skr, srs)
// }

// SkippedSummary ...
func (st *Stat) SkippedSummary() {

}

// PendingSummary ...
func (st *Stat) PendingSummary() {

}

// ScheduledSummary ...
func (st *Stat) ScheduledSummary() {

}
