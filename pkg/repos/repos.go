// Package repos collects Git repositories as Repo structs.
package repos

import (
	"bytes"
	"fmt"
	"log"
	"sort"
	"strconv"
	"sync"

	"github.com/jychri/git-in-sync/pkg/brf"
	"github.com/jychri/git-in-sync/pkg/conf"
	e "github.com/jychri/git-in-sync/pkg/emoji"
	"github.com/jychri/git-in-sync/pkg/flags"
	"github.com/jychri/git-in-sync/pkg/repo"
	"github.com/jychri/git-in-sync/pkg/stat"
	"github.com/jychri/git-in-sync/pkg/timer"
)

// private

func (rs Repos) names() (rss []string) {

	for _, r := range rs {
		rss = append(rss, r.Name)
	}

	if l := len(rss); l == 0 {
		log.Fatalf("No repos. Exiting")
	}

	return brf.Reduce(rss)
}

func (rs Repos) workspaces() (wss []string) {

	for _, r := range rs {
		wss = append(wss, r.Workspace)
	}

	if l := len(wss); l == 0 {
		log.Fatalf("No workspaces. Exiting")
	}

	return brf.Reduce(wss)
}

func (rs Repos) byName() {
	sort.SliceStable(rs, func(i, j int) bool { return rs[i].Name < rs[j].Name })
}

func (rs Repos) byWorkspacePath() {
	sort.SliceStable(rs, func(i, j int) bool { return rs[i].Name < rs[j].Name })
	sort.SliceStable(rs, func(i, j int) bool { return rs[i].WorkspacePath < rs[j].WorkspacePath })
}

func (rs Repos) syncVerifyWorkspaces(f flags.Flags, st *stat.Stat) {

	// sort Repos A-Z by *r.WorkspacePath
	rs.byWorkspacePath()

	// string slice collecting *r.Workspace(s)
	st.TotalWorkspaces = rs.workspaces()

	// "verifying workspaces ..."
	efc := e.Get("FileCabinet")
	l := len(st.TotalWorkspaces)
	sm := brf.Summary(st.TotalWorkspaces, 25)
	brf.Printv(f, "%v  verifying workspaces [%v](%v)", efc, l, sm)

	// verify each workspace, create if missing
	for _, r := range rs {
		r.VerifyWorkspace(f, st)
	}
}

func (rs Repos) summaryVerifyWorkspaces(f flags.Flags, st *stat.Stat, ti *timer.Timer) {
	vw := len(st.VerifiedWorkspaces)
	tw := len(st.TotalWorkspaces)
	cw := len(st.CreatedWorkspaces)

	// summary
	var b bytes.Buffer

	if vw == tw {
		b.WriteString(e.Get("Briefcase"))
	} else {
		b.WriteString(e.Get("Slash"))
	}

	b.WriteString(fmt.Sprintf(" [%v/%v] divs verified", vw, tw))

	if len(st.CreatedWorkspaces) >= 1 {
		b.WriteString(fmt.Sprintf(", created [%v]", strconv.Itoa(cw)))
	}

	b.WriteString(fmt.Sprintf(" {%v/%v}", ti.Split().String(), ti.Time().String()))

	brf.Printv(f, b.String())
}

func (rs Repos) syncVerifyRepos(f flags.Flags, st *stat.Stat) {
	for _, r := range rs {
		r.VerifyRepo(f, st)
	}
}

// Async

func (rs Repos) asyncClone(f flags.Flags, st *stat.Stat, t *timer.Timer) {

	if len(st.PendingClones) > 1 {
		es := e.Get("Sheep")
		pc := len(st.PendingClones)
		brf.Printv(f, "%v cloning [%v]", es, pc)

		var wg sync.WaitGroup
		for i := range rs {
			wg.Add(1)
			go func(r *repo.Repo) {
				defer wg.Done()
				r.GitClone(f, st)
			}(rs[i])
		}
		wg.Wait()

		t.Mark("async-clone")
	}
}

func (rs Repos) asyncInfo() {
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

// Public

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
	ws := rs.workspaces()

	// "workspaces|repos..."
	efm := e.Get("FaxMachine")
	lw := len(ws)
	lr := len(rs)
	ts := t.Split()
	tt := t.Time()
	brf.Printv(f, "%v [%v|%v] workspaces|repos {%v / %v}", efm, lw, lr, ts, tt)

	return rs
}

// VerifyWorkspaces verifies WorkspacePaths for Repos in Repos.
func (rs Repos) VerifyWorkspaces(f flags.Flags, st *stat.Stat, t *timer.Timer) {

	// verify each workspace
	rs.syncVerifyWorkspaces(f, st) // workspaceVerify

	// summary
	rs.summaryVerifyWorkspaces(f, st, t) // workspaceSummary
}

// VerifyRepos verifies all Repos in Repos.
func (rs Repos) VerifyRepos(f flags.Flags, st *stat.Stat, t *timer.Timer) {

	// check if present
	rs.syncVerifyRepos(f, st) // repoCheck

	// clone missing repos
	rs.asyncClone(f, st, t) // repoClone

	// get info for all repos
	rs.asyncInfo() // repoInfo

	// summary
	// rs.summaryVerifyRepos(f, st, t) -> create // repoSummary
}
