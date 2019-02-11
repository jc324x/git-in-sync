package stat

import (
	"bytes"
	"strconv"

	"github.com/jychri/git-in-sync/pkg/brf"
	e "github.com/jychri/git-in-sync/pkg/emoji"
	"github.com/jychri/git-in-sync/pkg/flags"
	"github.com/jychri/git-in-sync/pkg/timer"
)

// Stat holds values for the current run.
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

// Reduce reduces slices in *Run to their unique elements - no duplicates.
func (st *Stat) Reduce() {
	st.CreatedWorkspaces = brf.Reduce(st.CreatedWorkspaces)
	st.VerifiedWorkspaces = brf.Reduce(st.VerifiedWorkspaces)
	st.InaccessibleWorkspaces = brf.Reduce(st.InaccessibleWorkspaces)
	st.TotalWorkspaces = brf.Reduce(st.TotalWorkspaces)
	st.PendingClones = brf.Reduce(st.PendingClones)
	st.ClonedRepos = brf.Reduce(st.ClonedRepos)
}

// VCSummary ...
func (st *Stat) VCSummary(f flags.Flags, t *timer.Timer) {

	et := e.Get("Truck")
	cr := len(st.ClonedRepos)
	pc := len(st.PendingClones)
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
