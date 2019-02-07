// Package repos collects Git repositories as Repo structs.
package repos

import (
	"fmt"
	"log"
	"sort"
	"sync"

	"github.com/jychri/git-in-sync/pkg/brf"
	"github.com/jychri/git-in-sync/pkg/conf"
	"github.com/jychri/git-in-sync/pkg/e"
	"github.com/jychri/git-in-sync/pkg/flags"
	"github.com/jychri/git-in-sync/pkg/repo"
	"github.com/jychri/git-in-sync/pkg/run"
	"github.com/jychri/git-in-sync/pkg/timer"
)

// Repos collects pointers to Repo structs.
type Repos []*repo.Repo

// Init returns a slice of Repo structs.
func Init(c conf.Config, f flags.Flags, t *timer.Timer) (rs Repos) {

	ep := e.Get("Pager")
	brf.Printv(f, "%v parsing workspaces|repos", ep)

	// initialize Repos from Config
	for _, bl := range c.Bundles {
		for _, z := range bl.Zones {
			for _, rn := range z.Repos {
				r := repo.Init(z.Workspace, z.User, z.Remote, bl.Path, rn)
				rs = append(rs, r)
			}
		}
	}

	if l := len(rs); l == 0 {
		log.Fatalf("No repos. Exiting")
	}

	// timer
	t.Mark("init-repos")

	// sort
	ws := rs.Workspaces()

	// "workspaces|repos..."
	efm := e.Get("FaxMachine")
	lw := len(ws)
	lr := len(rs)
	ts := t.Split()
	tt := t.Time()
	brf.Printv(f, "%v [%v|%v] workspaces|repos {%v / %v}", efm, lw, lr, ts, tt)

	return rs
}

// Names returns all Repo Names.
func (rs Repos) Names() []string {
	var rss []string

	for _, r := range rs {
		rss = append(rss, r.Name)
	}

	if l := len(rss); l == 0 {
		log.Fatalf("No repos. Exiting")
	}

	return brf.Reduce(rss)
}

// Workspaces returns all Repo Workpaces.
func (rs Repos) Workspaces() []string {
	var wss []string

	for _, r := range rs {
		wss = append(wss, r.Workspace)
	}

	if l := len(wss); l == 0 {
		log.Fatalf("No workspaces. Exiting")
	}

	return brf.Reduce(wss)
}

// ByName sorts Repos in Repos A-Z by Name.
func (rs Repos) ByName() {
	sort.SliceStable(rs, func(i, j int) bool { return rs[i].Name < rs[j].Name })
}

// ByWorkspacePath sorts Repos in Repos A-Z by WorkspacePath.
func (rs Repos) ByWorkspacePath() {
	sort.SliceStable(rs, func(i, j int) bool { return rs[i].Name < rs[j].Name })
	sort.SliceStable(rs, func(i, j int) bool { return rs[i].WorkspacePath < rs[j].WorkspacePath })
}

// Async

// AsyncClone asynchronously clones all absent Repos in Repos.
func (rs Repos) AsyncClone(f flags.Flags, ru *run.Run, t *timer.Timer) {

	// "cloning ..."
	es := e.Get("Sheep")
	pc := len(ru.PendingClones)
	brf.Printv(f, "%v cloning [%v]", es, pc)

	var wg sync.WaitGroup
	for i := range rs {
		wg.Add(1)
		go func(r *repo.Repo) {
			defer wg.Done()
			r.GitClone(f, ru)
		}(rs[i])
	}
	wg.Wait()

	t.Mark("async-clone")
	ru.VCSummary(f, t)
}

// AsyncInfo asynchronously gathers info on all Repos in Repos.
func (rs Repos) AsyncInfo() {
	var wg sync.WaitGroup
	for i := range rs {
		wg.Add(1)
		go func(r *repo.Repo) {
			defer wg.Done()
			r.GitConfigOriginURL()
		}(rs[i])
	}
	wg.Wait()
}

// Main functions

// VerifyWorkspaces verifies WorkspacePaths for Repos in Repos.
func (rs Repos) VerifyWorkspaces(f flags.Flags, ru *run.Run, t *timer.Timer) {

	// sort Repos A-Z by *r.WorkspacePath
	rs.ByWorkspacePath()

	// []string of *r.Workspace
	ru.TotalWorkspaces = rs.Workspaces()

	ru.Reduce()

	// "verifying workspaces ..."
	efc := e.Get("FileCabinet")
	l := len(ru.TotalWorkspaces)
	sm := brf.Summary(ru.TotalWorkspaces, 25)
	brf.Printv(f, "%v  verifying workspaces [%v](%v)", efc, l, sm)

	for _, r := range rs {
		r.VerifyWorkspace(f, ru)
	}

	// print summary
	ru.VWSummary(f, t)
}

// VerifyRepos verifies all Repos in Repos.
func (rs Repos) VerifyRepos(f flags.Flags, ru *run.Run, t *timer.Timer) {

	for _, r := range rs {
		r.VerifyRepo(f, ru)
	}

	// async clone
	if len(ru.PendingClones) > 1 {
		rs.AsyncClone(f, ru, t)
	}

	// async info
	rs.AsyncInfo()
}
