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

func (rs Repos) repos() (rss []string) {

	for _, r := range rs {
		rss = append(rss, r.Name)
	}

	if l := len(rss); l == 0 {
		log.Fatalf("No repos. Exiting")
	}

	return rss
}

func initPrint(f flags.Flags) {
	ep := e.Get("Pager")
	brf.Printv(f, "%v parsing workspaces|repos", ep)
}

func initConvert(c conf.Config) (rs Repos) {
	for _, bl := range c.Bundles {
		for _, z := range bl.Zones {
			for _, rn := range z.Repos {
				r := repo.Init(z.Workspace, z.User, z.Remote, bl.Path, rn)
				rs = append(rs, r)
			}
		}
	}

	if len(rs) == 0 {
		log.Fatalf("No repos. Exiting")
	}

	return rs
}

func initSummary(f flags.Flags, st *stat.Stat, ti *timer.Timer, rs Repos) {
	efm := e.Get("FaxMachine")
	lw := len(st.Workspaces)
	lr := len(rs)
	ts := ti.Split()
	tt := ti.Time()
	brf.Printv(f, "%v [%v|%v] workspaces|repos {%v / %v}", efm, lw, lr, ts, tt)
}

func (rs Repos) byName() {
	sort.SliceStable(rs, func(i, j int) bool { return rs[i].Name < rs[j].Name })
}

func (rs Repos) byWorkspacePath() {
	sort.SliceStable(rs, func(i, j int) bool { return rs[i].Name < rs[j].Name })
	sort.SliceStable(rs, func(i, j int) bool { return rs[i].WorkspacePath < rs[j].WorkspacePath })
}

func (rs Repos) workspaceSync(f flags.Flags, st *stat.Stat, ti *timer.Timer) {

	// sort Repos A-Z by *r.WorkspacePath
	rs.byWorkspacePath()

	// "verifying workspaces ..."
	efc := e.Get("FileCabinet")
	l := len(st.Workspaces)
	sm := brf.Summary(st.Workspaces, 25)
	brf.Printv(f, "%v  verifying workspaces [%v](%v)", efc, l, sm)

	// verify each workspace, create if missing
	for _, r := range rs {
		r.VerifyWorkspace(f, st)
	}

	ti.Mark("workspace-sync")
}

func (rs Repos) workspaceSummary(f flags.Flags, st *stat.Stat, ti *timer.Timer) {
	vw := len(st.VerifiedWorkspaces)
	tw := len(st.Workspaces)
	cw := len(st.CreatedWorkspaces)

	// summary
	var b bytes.Buffer

	if vw == tw {
		b.WriteString(e.Get("Briefcase"))
	} else {
		b.WriteString(e.Get("Slash"))
	}

	b.WriteString(fmt.Sprintf(" [%v/%v] workspaces verified", vw, tw))

	if len(st.CreatedWorkspaces) >= 1 {
		b.WriteString(fmt.Sprintf(", created [%v]", strconv.Itoa(cw)))
	}

	b.WriteString(fmt.Sprintf(" {%v/%v}", ti.Split().String(), ti.Time().String()))

	brf.Printv(f, b.String())
}

func (rs Repos) cloneSchedule(f flags.Flags, st *stat.Stat) {
	for _, r := range rs {
		r.GitSchedule(f, st)
	}
}

// Async

func (rs Repos) cloneAsync(f flags.Flags, st *stat.Stat, ti *timer.Timer) {

	if len(st.PendingClones) > 1 {
		es := e.Get("Sheep")
		pc := len(st.PendingClones)
		brf.Printv(f, "%v cloning [%v]", es, pc)

		var wg sync.WaitGroup
		for i := range rs {
			wg.Add(1)
			go func(r *repo.Repo) {
				defer wg.Done()
				r.GitClone(f)
			}(rs[i])
		}
		wg.Wait()
	}

	ti.Mark("async-clone")
}

func (rs Repos) cloneSummary(f flags.Flags, st *stat.Stat, ti *timer.Timer) {

	for _, r := range rs {
		if r.Cloned == true {
			st.ClonedRepos = append(st.ClonedRepos, r.Name)
		}
	}

	if len(st.ClonedRepos) == 0 {
		return
	}

	et := e.Get("Truck")
	lc := len(st.ClonedRepos)
	lp := len(st.PendingClones)
	ts := ti.Split()
	tt := ti.Time()
	brf.Printv(f, "%v [%v/%v] repos cloned {%v / %v}", et, lc, lp, ts, tt)
}

func (rs Repos) infoAsync(f flags.Flags) {
	var wg sync.WaitGroup
	for i := range rs {
		wg.Add(1)
		go func(r *repo.Repo) {
			defer wg.Done()
			r.GitConfigOriginURL()
			r.GitRemoteUpdate()
			r.GitAbbrevRef()
			r.GitLocalSHA()
			r.GitUpstreamBranch()
			r.GitMergeBaseSHA()
			r.GitRevParseUpstream()
			r.GitDiffsNameOnly()
			r.GitShortstat()
			r.GitUntracked()
			r.SetStatus(f)
		}(rs[i])
	}
	wg.Wait()
}

func (rs Repos) repoSummary(f flags.Flags, st *stat.Stat, ti *timer.Timer) {
	st.Repos = rs.repos()

	for _, r := range rs {
		switch r.Category {
		case "Pending":
			st.PendingRepos = append(st.PendingRepos, r.Name)
		case "Skipped":
			st.SkippedRepos = append(st.SkippedRepos, r.Name)
		case "Scheduled":
			st.ScheduledRepos = append(st.ScheduledRepos, r.Name)
		case "Complete":
			st.CompleteRepos = append(st.CompleteRepos, r.Name)
		}
	}

	tr := len(st.Repos)
	cr := len(st.CompleteRepos)

	if tr == cr {
		st.Complete = true
		ec := e.Get("Checkmark")
		brf.Printv(f, "%v [%v/%v] complete", ec, cr, tr)
		return
	}

	st.Complete = false

	var b bytes.Buffer
	pr := len(st.PendingRepos)
	skr := len(st.SkippedRepos)
	scr := len(st.ScheduledRepos)
	ew := e.Get("Warning")

	b.WriteString(ew)
	b.WriteString(" [")
	b.WriteString(strconv.Itoa(cr))
	b.WriteString(" / ")
	b.WriteString(strconv.Itoa(tr))
	b.WriteString("] ")

	if pr >= 1 {

	}

	if skr >= 1 {

	}

	if scr >= 1 {

	}

}

// Public

// Repos collects pointers to Repo structs.
type Repos []*repo.Repo

// Init returns a slice of Repo structs.
func Init(c conf.Config, f flags.Flags, st *stat.Stat, ti *timer.Timer) Repos {
	initPrint(f)                    // print startup
	rs := initConvert(c)            // convert Config into Repos
	st.Workspaces = rs.workspaces() // record stats
	ti.Mark("init-repos")           // mark timer
	initSummary(f, st, ti, rs)      // print summary
	return rs
}

// VerifyWorkspaces verifies WorkspacePaths for Repos in Repos.
func (rs Repos) VerifyWorkspaces(f flags.Flags, st *stat.Stat, ti *timer.Timer) {
	rs.workspaceSync(f, st, ti)    // create missing workspaces
	rs.workspaceSummary(f, st, ti) // print summary
}

// VerifyRepos verifies all Repos in Repos.
func (rs Repos) VerifyRepos(f flags.Flags, st *stat.Stat, ti *timer.Timer) {
	rs.cloneSchedule(f, st)    // mark pending clones
	rs.cloneAsync(f, st, ti)   // clone missing repos (async)
	rs.cloneSummary(f, st, ti) // print summary
	rs.infoAsync(f)            // get info for all repos (async)
	rs.repoSummary(f, st, ti)  // print summary
}
