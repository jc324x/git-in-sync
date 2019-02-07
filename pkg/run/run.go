package run

import (
	"bytes"
	// "fmt"
	"strconv"

	"github.com/jychri/git-in-sync/pkg/brf"
	"github.com/jychri/git-in-sync/pkg/e"
	"github.com/jychri/git-in-sync/pkg/flags"
	"github.com/jychri/git-in-sync/pkg/timer"
)

// Run holds values for the current run.
type Run struct {
	CreatedWorkspaces      []string
	VerifiedWorkspaces     []string
	InaccessibleWorkspaces []string
	TotalWorkspaces        []string
	PendingClones          []string
	ClonedRepos            []string
	TotalRepos             []string
}

// Init returns a new *Run.
func Init() *Run {
	ru := new(Run)
	return ru
}

// Reduce reduces slices in *Run to their unique elements - no duplicates.
func (ru *Run) Reduce() {
	ru.CreatedWorkspaces = brf.Reduce(ru.CreatedWorkspaces)
	ru.VerifiedWorkspaces = brf.Reduce(ru.VerifiedWorkspaces)
	ru.InaccessibleWorkspaces = brf.Reduce(ru.InaccessibleWorkspaces)
	ru.TotalWorkspaces = brf.Reduce(ru.TotalWorkspaces)
	ru.PendingClones = brf.Reduce(ru.PendingClones)
	ru.ClonedRepos = brf.Reduce(ru.ClonedRepos)
}

// VCSummary ...
func (ru *Run) VCSummary(f flags.Flags, t *timer.Timer) {

	et := e.Get("Truck")
	cr := len(ru.ClonedRepos)
	pc := len(ru.PendingClones)
	ts := t.Split().Truncate(timer.M).String()
	tt := t.Time().Truncate(timer.M).String()

	var b bytes.Buffer
	b.WriteString(et)
	b.WriteString(" [")
	b.WriteString(strconv.Itoa(cr))
	b.WriteString("/")
	b.WriteString(strconv.Itoa(pc))
	b.WriteString("] cloned")

	b.WriteString(" {")
	b.WriteString(ts)
	b.WriteString(" / ")
	b.WriteString(tt)
	b.WriteString("}")

	brf.Printv(f, b.String())
}
