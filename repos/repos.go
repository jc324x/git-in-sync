// Package repos collects Git repositories as Repo structs.
package repos

import (
	"bytes"
	"fmt"
	"log"
	"sort"
	"sync"

	"github.com/jychri/git-in-sync/brf"
	"github.com/jychri/git-in-sync/conf"
	"github.com/jychri/git-in-sync/emoji"
	"github.com/jychri/git-in-sync/flags"
	"github.com/jychri/git-in-sync/repo"
	"github.com/jychri/git-in-sync/stat"
	"github.com/jychri/git-in-sync/timer"
)

// private

// direct returns a *Repo from Repos with a matching name
func (rs Repos) direct(name string) *repo.Repo {
	for _, r := range rs {
		if name == r.Name {
			return r
		}
	}
	return nil
}

// names returns a string slice of all Repo names.
func (rs Repos) names() (rss []string) {

	for _, r := range rs {
		rss = append(rss, r.Name)
	}

	if l := len(rss); l == 0 {
		log.Fatalf("No repos. Exiting")
	}

	return brf.Reduce(rss)
}

// names returns a string slice of all Repo workspaces.
func (rs Repos) workspaces() (wss []string) {

	for _, r := range rs {
		wss = append(wss, r.Workspace)
	}

	if l := len(wss); l == 0 {
		log.Fatalf("No workspaces. Exiting")
	}

	return brf.Reduce(wss)
}

// sort repos by Name
func (rs Repos) byName() {
	sort.SliceStable(rs, func(i, j int) bool { return rs[i].Name < rs[j].Name })
}

// sort repos by WorkspacePath
func (rs Repos) byWorkspacePath() {
	sort.SliceStable(rs, func(i, j int) bool { return rs[i].Name < rs[j].Name })
	sort.SliceStable(rs, func(i, j int) bool { return rs[i].WorkspacePath < rs[j].WorkspacePath })
}

// print startup
func initPrint(f flags.Flags) {
	ep := emoji.Get("Pager") // Pager emoji
	flags.Printv(f, "%v parsing workspaces|repos", ep)
}

// convert Config to Repos
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

// print summary
func initSummary(f flags.Flags, st *stat.Stat, ti *timer.Timer, rs Repos) {
	efm := emoji.Get("FaxMachine") // FaxMachine emoji
	lw := len(st.Workspaces)       // number of workspaces
	lr := len(rs)                  // number of repos
	ts := ti.Split()               // last split
	tt := ti.Time()                // elapsed time
	flags.Printv(f, "%v [%v|%v] workspaces|repos {%v / %v}", efm, lw, lr, ts, tt)
}

// create missing workspaces
func (rs Repos) workspaceSync(f flags.Flags, st *stat.Stat, ti *timer.Timer) {

	rs.byWorkspacePath()                 // sort Repos A-Z by WorkspacePath
	efc := emoji.Get("FileCabinet")      // FileCabinet emoji
	tw := len(st.Workspaces)             // number of workspaces
	sm := brf.Summary(st.Workspaces, 25) // summary of workspaces
	flags.Printv(f, "%v  verifying workspaces [%v](%v)", efc, tw, sm)

	// verify each workspace, create if missing
	for _, r := range rs {
		r.VerifyWorkspace(f, st)
	}

	ti.Mark("workspace-sync") // mark workspace-sync
}

// print summary
func (rs Repos) workspaceSummary(f flags.Flags, st *stat.Stat, ti *timer.Timer) {

	st.Reduce()                      // reduce slices to their unique items
	tw := len(st.Workspaces)         // number of workspaces
	vw := len(st.VerifiedWorkspaces) // number of verified workspaces
	cw := len(st.CreatedWorkspaces)  // number of created workspaces
	eb := emoji.Get("Briefcase")     // Briefcase emoji
	es := emoji.Get("Slash")         // Slash emoji
	ts := ti.Split()                 // last split
	tt := ti.Time()                  // elapsed time

	var b bytes.Buffer // buffer

	// Briefcase for match, Slash for mismatch
	if vw == tw {
		b.WriteString(eb)
	} else {
		b.WriteString(es)
	}

	b.WriteString(fmt.Sprintf(" [%v/%v] workspaces verified", vw, tw))

	// only print ", created ..." if needed
	if cw >= 1 {
		b.WriteString(fmt.Sprintf(", created [%v]", cw))
	}

	// write timer info
	b.WriteString(fmt.Sprintf(" {%v/%v}", ts, tt))

	flags.Printv(f, b.String())
}

// schedule pending clones
func (rs Repos) cloneSchedule(f flags.Flags, st *stat.Stat) {
	for _, r := range rs {
		r.GitSchedule(f, st)
	}
}

// clone missing repos (async)
func (rs Repos) cloneAsync(f flags.Flags, st *stat.Stat, ti *timer.Timer) {

	// return early if no pending clones
	if len(st.PendingClones) == 0 {
		return
	}

	es := emoji.Get("Sheep")                // Sheep emoji
	pc := len(st.PendingClones)             // number of pending clones
	ps := brf.Summary(st.PendingClones, 25) // short summary
	flags.Printv(f, "%v cloning [%v](%v)", es, pc, ps)

	var wg sync.WaitGroup
	for i := range rs {
		wg.Add(1)
		go func(r *repo.Repo) {
			defer wg.Done()
			r.GitClone(f)
		}(rs[i])
	}
	wg.Wait()

	ti.Mark("async-clone") // mark async-clone
}

// print summary
func (rs Repos) cloneSummary(f flags.Flags, st *stat.Stat, ti *timer.Timer) {

	// loop over repos again, record clone count in Stat.
	for _, r := range rs {
		if r.Cloned == true {
			st.ClonedRepos = append(st.ClonedRepos, r.Name)
		}
	}

	// return if nothing was cloned
	if len(st.ClonedRepos) == 0 {
		return
	}

	et := emoji.Get("Truck")    // Truck emoji
	lc := len(st.ClonedRepos)   // number of cloned repos
	lp := len(st.PendingClones) // number of pending clones
	ts := ti.Split()            // last split
	tt := ti.Time()             // elapsed time
	flags.Printv(f, "%v [%v/%v] repos cloned {%v / %v}", et, lc, lp, ts, tt)
}

// print startup
func (rs Repos) infoPrint(f flags.Flags, st *stat.Stat) {
	st.Repos = rs.names()           // record repo names
	ep := emoji.Get("Satellite")    // Satellite emoji
	lr := len(st.Repos)             // number of repos
	sr := brf.Summary(st.Repos, 25) // short summary
	flags.Printv(f, "%v  checking repos [%v](%v)", ep, lr, sr)
}

func (rs Repos) infoAsync(f flags.Flags, ti *timer.Timer) {

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

	ti.Mark("info-async") // mark info-async
}

func (rs Repos) changesSummary(f flags.Flags, st *stat.Stat, ti *timer.Timer) {
	rs.category(st)              // set current category stats
	tr := len(st.Repos)          // number of repos
	cr := len(st.CompleteRepos)  // number of complete repos
	ec := emoji.Get("Checkmark") // Checkmark emoji
	es := emoji.Get("Fire")      // Fire emoji (replace with stop sign)
	ew := emoji.Get("Warning")   // Warning emoji
	ti.Mark("repo-summary")      // mark repo-summary
	ts := ti.Split()             // last split
	tt := ti.Time()              // elapsed time

	switch {
	case st.CheckComplete():
		flags.Printv(f, "%v [%v/%v] repos complete {%v / %v}", ec, cr, tr, ts, tt)
	case st.CheckPending():
		flags.Printv(f, "%v [%v/%v] repos complete {%v / %v}", ew, cr, tr, ts, tt)
	case st.CheckSkipped():
		flags.Printv(f, "%v [%v/%v] repos complete {%v / %v}", es, cr, tr, ts, tt)
	}

}

// infoSummary ...
func (rs Repos) infoSummary(f flags.Flags, st *stat.Stat, ti *timer.Timer) {
	rs.category(st)              // set current category stats
	tr := len(st.Repos)          // number of repos
	cr := len(st.CompleteRepos)  // number of complete repos
	ec := emoji.Get("Checkmark") // Checkmark emoji
	es := emoji.Get("Fire")      // Fire emoji (replace with stop sign)
	ew := emoji.Get("Warning")   // Warning emoji
	ti.Mark("repo-summary")      // mark repo-summary
	ts := ti.Split()             // last split
	tt := ti.Time()              // elapsed time

	switch {
	case st.CheckComplete():
		flags.Printv(f, "%v [%v/%v] repos complete {%v / %v}", ec, cr, tr, ts, tt)
	case st.CheckPending():
		flags.Printv(f, "%v [%v/%v] repos complete {%v / %v}", ew, cr, tr, ts, tt)
	case st.CheckSkipped():
		flags.Printv(f, "%v [%v/%v] repos complete {%v / %v}", es, cr, tr, ts, tt)
	}
}

func (rs Repos) promptUser(f flags.Flags, st *stat.Stat) {
	if st.CheckComplete() {
		return
	}

	for _, r := range rs {
		r.UserConfirm(f)
	}
}

func (rs Repos) changesAsync(f flags.Flags, st *stat.Stat, ti *timer.Timer) {
	if st.CheckComplete() {
		return
	}

	var wg sync.WaitGroup
	for i := range rs {
		wg.Add(1)
		go func(r *repo.Repo) {
			defer wg.Done()

			if r.Category != "Scheduled" {
				return
			}

			switch r.Action {
			case "Pull":
				r.GitPull(f)
				r.GitClear()
			case "Push":
				r.GitPush(f)
				r.GitClear()
			case "Add-Commit-Push":
				r.GitAdd(f)
				r.GitCommit(f)
				r.GitPush(f)
				r.GitClear()
			case "Stash-Pull-Pop-Commit-Push":
				r.GitAdd(f)
				r.GitStash(f)
				r.GitPull(f)
				r.GitPop(f)
				r.GitAdd(f)
				r.GitCommit(f)
				r.GitPush(f)
				r.GitClear()
			}
		}(rs[i])
	}
	wg.Wait()

	st.Clear()
}

// update stat.Stat, collect Repos by category
func (rs Repos) category(st *stat.Stat) {
	for _, r := range rs {
		switch {
		case r.Category == "Pending":
			st.PendingRepos = append(st.PendingRepos, r.Name)
		case r.Category == "Skipped":
			st.SkippedRepos = append(st.SkippedRepos, r.Name)
		case r.Category == "Scheduled" && r.Action == "Push":
			st.ScheduledPush = append(st.ScheduledRepos, r.Name)
		case r.Category == "Scheduled" && r.Action == "Pull":
			st.ScheduledPull = append(st.ScheduledRepos, r.Name)
		case r.Category == "Complete":
			st.CompleteRepos = append(st.CompleteRepos, r.Name)
		}
	}
}

// Public

// Repos collects pointers to Repo structs.
type Repos []*repo.Repo

// Init returns a slice of Repo structs.
func Init(c conf.Config, f flags.Flags, st *stat.Stat, ti *timer.Timer) Repos {
	initPrint(f)                    // print startup
	rs := initConvert(c)            // convert Config to Repos
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
	rs.cloneSchedule(f, st)    // schedule pending clones
	rs.cloneAsync(f, st, ti)   // clone missing repos (async)
	rs.cloneSummary(f, st, ti) // print summary
	rs.infoPrint(f, st)        // print startup
	rs.infoAsync(f, ti)        // update info (async)
	rs.infoSummary(f, st, ti)  // print summary
}

// VerifyChanges ...
func (rs Repos) VerifyChanges(f flags.Flags, st *stat.Stat, ti *timer.Timer) {
	rs.promptUser(f, st)         // prompt user
	rs.changesAsync(f, st, ti)   // submit changes (async)
	rs.infoAsync(f, ti)          // update info (async)
	rs.changesSummary(f, st, ti) // update info (async)
}
