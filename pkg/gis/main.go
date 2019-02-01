package main

import (
	"github.com/jychri/git-in-sync/pkg/brf"
	"github.com/jychri/git-in-sync/pkg/conf"
	e "github.com/jychri/git-in-sync/pkg/emoji"
	"github.com/jychri/git-in-sync/pkg/flags"
	"github.com/jychri/git-in-sync/pkg/repos"
	"github.com/jychri/git-in-sync/pkg/timer"
)

// Init ...
func Init() (f flags.Flags, rs repos.Repos, t *timer.Timer) {

	e.ClearScreen()

	// initialize Timer and Flags
	t = timer.Init()
	f = flags.Init()
	t.MarkMoment("init-flags")

	// "start"
	brf.Printv(f, "%v start", e.Get("Clapper"))

	// "flag(s) set..."
	if ft, err := t.GetMoment("init-flags"); err == nil {
		brf.Printv(f, "%v parsing flags", e.Get("FlagInHole"))
		brf.Printv(f, "%v running in '%v' mode {%v / %v}", e.Get("Flag"), f.Mode, ft.Split, ft.Start)
	}

	// "reading ~/.gisrc.json"
	brf.Printv(f, "%v reading ~/.gisrc.json", e.Get("Books"))

	// initialize Config from conf.Path(f)
	c := conf.Init(f)
	t.MarkMoment("init-config")

	// "read conf.Path(f)"
	brf.Printv(f, "%v read %v {%v / %v}", e.Get("Book"), conf.Path(f), t.Split(), t.Time())

	// initialize Repos
	rs = repos.Init(c)

	return f, rs, t
}

func main() {
	f, rs, t := Init()
	rs.VerifyDivs(f, t)
	// rs.VerifyCloned(e, f, t)
	// rs.VerifyRepos(e, f, t)
	// rs.VerifyChanges(e, f, t)
	// rs.SubmitChanges(e, f, t)
}
