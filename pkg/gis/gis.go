// Package gis is the main package of git-in-sync.
package main

import (
	"github.com/jychri/git-in-sync/pkg/brf"
	"github.com/jychri/git-in-sync/pkg/conf"
	"github.com/jychri/git-in-sync/pkg/emoji"
	"github.com/jychri/git-in-sync/pkg/flags"
	"github.com/jychri/git-in-sync/pkg/repos"
	"github.com/jychri/git-in-sync/pkg/stat"
	"github.com/jychri/git-in-sync/pkg/timer"
)

// Init initializes Flags, Repos, Stat and a Timer.
func Init() (f flags.Flags, rs repos.Repos, st *stat.Stat, ti *timer.Timer) {

	emoji.ClearScreen()                                                // clear screen
	ti = timer.Init()                                                  // init Timer
	f = flags.Init()                                                   // init Flags
	ti.Mark("init-flags")                                              // mark init flags
	st = stat.Init()                                                   // init Stat
	ec := emoji.Get("Clapper")                                         // print start
	brf.Printv(f, "%v start", ec)                                      //
	efih := emoji.Get("FlagInHole")                                    // print flags
	brf.Printv(f, "%v parsing flags", efih)                            //
	ef := emoji.Get("Flag")                                            //
	fm := f.Mode                                                       //
	ts := ti.Split()                                                   //
	tt := ti.Time()                                                    //
	brf.Printv(f, "%v running in '%v' mode {%v / %v}", ef, fm, ts, tt) //
	eb := emoji.Get("Books")                                           // print config
	brf.Printv(f, "%v reading ~/.gisrc.json", eb)                      //
	c := conf.Init(f)                                                  // init config
	ti.Mark("init-config")                                             // mark init config
	ebs := emoji.Get("Book")                                           //
	fc := f.Config                                                     //
	ts = ti.Split()                                                    //
	tt = ti.Time()                                                     //
	brf.Printv(f, "%v read %v {%v / %v}", ebs, fc, ts, tt)             //
	rs = repos.Init(c, f, st, ti)                                      // init repos
	return f, rs, st, ti                                               // return
}

func main() {
	f, rs, st, t := Init()
	rs.VerifyWorkspaces(f, st, t)
	rs.VerifyRepos(f, st, t)
	rs.VerifyChanges(f, st, t)
	rs.SubmitChanges(f, st, t)
}
