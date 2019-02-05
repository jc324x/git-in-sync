package run

import (
	"github.com/jychri/git-in-sync/pkg/brf"
)

// Run tracks stats for the current run
type Run struct {
	CWS []string // created workspaces
	VWS []string // verified workspaces
	IWS []string // inaccessible workspaces
	CWC int      // len(CWS)
	VWC int      // len(VWS)
	IWC int      // len(IWS)
	PCS []string // pending clones
	PCC int      // lens(PCS)
	CRS []string // cloned repos
	CRC int      // len(CRS)
}

// Init ...
func Init() *Run {
	ru := new(Run)
	return ru
}

// Reduce reduces slices in the *Run to their unique elements; no duplicates.
func (ru *Run) Reduce() {
	ru.CWS = brf.Reduce(ru.CWS)
	ru.VWS = brf.Reduce(ru.VWS)
	ru.IWS = brf.Reduce(ru.IWS)
	ru.CWC = len(ru.CWS)
	ru.VWC = len(ru.VWS)
	ru.IWC = len(ru.IWS)
	ru.PCC = len(ru.PCS)
	ru.CRC = len(ru.CRS)
}
