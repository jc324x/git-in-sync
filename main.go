// Package gis is the main package of git-in-sync.
package main

import (
	"github.com/jychri/git-in-sync/conf"
	"github.com/jychri/git-in-sync/emoji"
	"github.com/jychri/git-in-sync/flags"
	"github.com/jychri/git-in-sync/repos"
	"github.com/jychri/git-in-sync/stat"
	"github.com/jychri/git-in-sync/timer"
)

// Init initializes Flags, Repos, Stat and Timer.
func Init() (f flags.Flags, rs repos.Repos, st *stat.Stat, ti *timer.Timer) {
	ti = timer.Init()                                                    // init Timer
	f = flags.Init()                                                     // init Flags
	f.ClearScreen()                                                      // clear screen
	ti.Mark("init-flags")                                                // mark init flags
	st = stat.Init()                                                     // init Stat
	ec := emoji.Get("Clapper")                                           // print start
	flags.Printv(f, "%v start", ec)                                      // start message
	efih := emoji.Get("FlagInHole")                                      // FlagInHole emoji
	flags.Printv(f, "%v parsing flags", efih)                            // print parsing flags
	ef := emoji.Get("Flag")                                              // Flag emoji
	fm := f.Mode                                                         // flag.Mode
	ts := ti.Split()                                                     // last split
	tt := ti.Time()                                                      // elapsed time
	flags.Printv(f, "%v running in '%v' mode {%v / %v}", ef, fm, ts, tt) // print mode
	eb := emoji.Get("Books")                                             // Books emoji
	flags.Printv(f, "%v reading ~/.gisrc.json", eb)                      // print reading config
	c := conf.Init(f)                                                    // init config
	ti.Mark("init-config")                                               // mark init config
	ebs := emoji.Get("Book")                                             // Book emoji
	fc := f.Config                                                       // flag.Config
	ts = ti.Split()                                                      // last split
	tt = ti.Time()                                                       // elapsed time
	flags.Printv(f, "%v read %v {%v / %v}", ebs, fc, ts, tt)             // print read config
	rs = repos.Init(c, f, st, ti)                                        // init repos
	return f, rs, st, ti                                                 // return
}

func main() {
	f, rs, st, t := Init()        // init Flags, Repos, Stat and Timer
	rs.VerifyWorkspaces(f, st, t) // verify workspaces, create if needed
	rs.VerifyRepos(f, st, t)      // verify repos, clone if needed (async)
	rs.VerifyChanges(f, st, t)    // verify and submit changes (async)
}
