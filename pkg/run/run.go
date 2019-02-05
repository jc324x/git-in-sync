package run

import (
	"github.com/jychri/git-in-sync/pkg/brf"
)

// Run tracks stats for the current run
type Run struct {
	CreatedW      []string // created workspaces
	VerifiedW     []string // verified workspaces
	InaccessibleW []string // inaccessible workspaces
}

// Init ...
func Init() *Run {
	ru := new(Run)
	return ru
}

// Reduce reduces slices in the *Run to their unique elements; no duplicates.
func (ru *Run) Reduce() {
	ru.CreatedW = brf.Single(ru.CreatedW)
	ru.VerifiedW = brf.Single(ru.VerifiedW)
	ru.InaccessibleW = brf.Single(ru.InaccessibleW)
}
