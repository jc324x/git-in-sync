package run

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/jychri/git-in-sync/pkg/brf"
	"github.com/jychri/git-in-sync/pkg/e"
	"github.com/jychri/git-in-sync/pkg/flags"
	"github.com/jychri/git-in-sync/pkg/timer"
)

// Run tracks stats for the current run
type Run struct {
	CWS []string // created workspaces
	VWS []string // verified workspaces
	IWS []string // inaccessible workspaces
	TWS []string // total workspaces
	CWC int      // len(CWS)
	VWC int      // len(VWS)
	IWC int      // len(IWS)
	TWC int      // len(TWS)
	PCS []string // pending clones
	PCC int      // lens(PCS)
	CRS []string // cloned repos
	CRC int      // len(CRS)
}

// Init returns a new *Run.
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
	ru.TWC = len(ru.TWS)
	ru.PCC = len(ru.PCS)
	ru.CRC = len(ru.CRS)
}

// VWSummary ...
func (ru *Run) VWSummary(f flags.Flags, t *timer.Timer) {

	// summary
	var b bytes.Buffer

	if ru.TWC == ru.VWC {
		b.WriteString(e.Get("Briefcase"))
	} else {
		b.WriteString(e.Get("Slash"))
	}

	b.WriteString(fmt.Sprintf(" [%v/%v] divs verified", ru.VWC, ru.TWC))

	if ru.CWC >= 1 {
		b.WriteString(fmt.Sprintf(", created [%v]", strconv.Itoa(ru.CWC)))
	}

	b.WriteString(fmt.Sprintf(" {%v/%v}", t.Split().String(), t.Time().String()))

	brf.Printv(f, b.String())
}

// VCSummary ...
func (ru *Run) VCSummary(f flags.Flags, t *timer.Timer) {
	var b bytes.Buffer

	b.WriteString(e.Get("Truck"))
	b.WriteString(" [")
	b.WriteString(strconv.Itoa(ru.CRC))
	b.WriteString("/")
	b.WriteString(strconv.Itoa(ru.PCC))
	b.WriteString("] cloned")

	b.WriteString(" {")
	b.WriteString(t.Split().Truncate(timer.M).String())
	b.WriteString(" / ")
	b.WriteString(t.Time().Truncate(timer.M).String())
	b.WriteString("}")

	brf.Printv(f, b.String())
}
