// Gis is the nexus...
package main

import (
	"github.com/jychri/git-in-sync/pkg/brf"
	"github.com/jychri/git-in-sync/pkg/conf"
	e "github.com/jychri/git-in-sync/pkg/emoji"
	"github.com/jychri/git-in-sync/pkg/flags"
	"github.com/jychri/git-in-sync/pkg/repos"
	"github.com/jychri/git-in-sync/pkg/run"
	"github.com/jychri/git-in-sync/pkg/timer"
)

// Init returns initial
func Init() (f flags.Flags, rs repos.Repos, ru *run.Run, ti *timer.Timer) {

	e.ClearScreen()

	// initialize Timer and Flags
	ti = timer.Init()
	f = flags.Init()
	ru = run.Init()

	ti.Mark("init-flags")

	// "start"
	brf.Printv(f, "%v start", e.Get("Clapper"))

	// "flag(s) set..."
	if ft, err := ti.Get("init-flags"); err == nil {
		brf.Printv(f, "%v parsing flags", e.Get("FlagInHole"))
		brf.Printv(f, "%v running in '%v' mode {%v / %v}", e.Get("Flag"), f.Mode, ft.Split, ft.Start)
	}

	// "reading ~/.gisrc.json"
	brf.Printv(f, "%v reading ~/.gisrc.json", e.Get("Books"))

	// initialize Config from conf.Path(f)
	c := conf.Init(f)
	ti.Mark("init-config")

	// "read conf.Path(f)"
	brf.Printv(f, "%v read %v {%v / %v}", e.Get("Book"), (f.Config), ti.Split(), ti.Time())

	// initialize Repos
	rs = repos.Init(c, f, ti)

	return f, rs, ru, ti
}

func main() {
	f, rs, ru, t := Init()
	rs.VerifyWorkspaces(f, ru, t)
	rs.VerifyRepos(f, ru, t)
	// rs.VerifyChanges(f, ru, t)
	// rs.SubmitChanges(e, f, t)
}
