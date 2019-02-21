// Package repos collects Git repositories as Repo structs.
package repos

import (
	"bytes"
	"fmt"
	"log"
	"sort"
	"sync"

	"github.com/jychri/git-in-sync/pkg/brf"
	"github.com/jychri/git-in-sync/pkg/conf"
	"github.com/jychri/git-in-sync/pkg/emoji"
	"github.com/jychri/git-in-sync/pkg/flags"
	"github.com/jychri/git-in-sync/pkg/repo"
	"github.com/jychri/git-in-sync/pkg/stat"
	"github.com/jychri/git-in-sync/pkg/timer"
)

// private

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
	ep := emoji.Get("Pager") // Pager emoji
	flags.Printv(f, "%v parsing workspaces|repos", ep)
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
	efm := emoji.Get("FaxMachine") // FaxMachine emoji
	lw := len(st.Workspaces)       // number of workspaces
	lr := len(rs)                  // number of repos
	ts := ti.Split()               // last split
	tt := ti.Time()                // elapsed time
	flags.Printv(f, "%v [%v|%v] workspaces|repos {%v / %v}", efm, lw, lr, ts, tt)
}

func (rs Repos) byName() {
	sort.SliceStable(rs, func(i, j int) bool { return rs[i].Name < rs[j].Name })
}

func (rs Repos) byWorkspacePath() {
	sort.SliceStable(rs, func(i, j int) bool { return rs[i].Name < rs[j].Name })
	sort.SliceStable(rs, func(i, j int) bool { return rs[i].WorkspacePath < rs[j].WorkspacePath })
}

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

func (rs Repos) cloneSchedule(f flags.Flags, st *stat.Stat) {
	for _, r := range rs {
		r.GitSchedule(f, st)
	}
}

func (rs Repos) cloneAsync(f flags.Flags, st *stat.Stat, ti *timer.Timer) {

	if len(st.PendingClones) > 1 {
		es := emoji.Get("Sheep")                // Sheep emoji
		pc := len(st.PendingClones)             // number of pending clones
		ps := brf.Summary(st.PendingClones, 25) // short summary
		flags.Printv(f, "%v cloning [%v](%v)", es, pc, ps)

		// asynchorously clone missing repos. GitClone
		// returns early if PendingClone == false.
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

	ti.Mark("async-clone") // mark async-clone
}

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

func (rs Repos) infoPrint(f flags.Flags, st *stat.Stat) {
	st.Repos = rs.repos()           // record repo names
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

func (rs Repos) infoSummary(f flags.Flags, st *stat.Stat, ti *timer.Timer) {

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

	tr := len(st.Repos)          // number of repos
	cr := len(st.CompleteRepos)  // number of complete repos
	ec := emoji.Get("Checkmark") // Checkmark emoji
	es := emoji.Get("Fire")      // Fire emoji (replace with stop sign)
	ew := emoji.Get("Warning")   // Warning emoji
	ti.Mark("repo-summary")      // mark repo-summary
	ts := ti.Split()             // last split
	tt := ti.Time()              // elapsed time

	var b bytes.Buffer // buffer

	// Checkmark for complete, Warning for incomplete, Stop for skipped only
	switch {
	case st.AllComplete():
		b.WriteString(ec)
	case st.OnlySkipped():
		b.WriteString(es)
	default:
		b.WriteString(ew)
	}

	b.WriteString(fmt.Sprintf(" [%v/%v] repos complete {%v / %v}", cr, tr, ts, tt))
	flags.Printv(f, b.String())

	switch {
	case st.AllComplete():
		return
	case st.OnlySkipped():
		st.SkippedSummary()
	case st.OnlyScheduled():
		st.ScheduledSummary()
	}
}

func (rs Repos) promptUser(f flags.Flags) {
	for _, r := range rs {
		r.UserConfirm(f)
	}
}

func (rs Repos) submitChanges(f flags.Flags, st *stat.Stat, ti *timer.Timer) {
	if st.Continue() == false {
		return
	}

	var wg sync.WaitGroup
	for i := range rs {
		wg.Add(1)
		go func(r *repo.Repo) {
			defer wg.Done()
			switch r.Action {
			case "Pull":
				r.GitPull(f)
			case "Push":
				r.GitPush(f)
			case "Add-Commit-Push":
				r.GitAdd(f)
				r.GitCommit(f)
				r.GitPush(f)
			case "Stash-Pull-Pop-Commit-Push":
				r.GitAdd(f)
				r.GitStash(f)
				r.GitPull(f)
				r.GitPop(f)
				r.GitCommit(f)
				r.GitPush(f)
			}
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

	// var vc []string // verified complete repos

	// for _, r := range srs {
	// 	if r.Category == "Complete" {
	// 		vc = append(vc, r.Name)
	// 	}
	// }

	//
	// switch {
	// case len(srs) == len(vc) && len(skrs) == 0:
	// 	fmt.Println("all good. nothing skipped, everything completed")
	// // case len(srs) == len(vc) && len(skrs) >= 1:
	// // 	fmt.Println("all pending actions complete - did skip this though (as planned)")
	// case len(srs) != len(vc) && len(skrs) >= 1:
	// 	fmt.Println("all changes not submitted correctly, also skipped")
	// }

	// if len(srs) == len(vc) {
	// 	fmt.Println("All changes submitted for pending repos")
	// } else {
	// 	fmt.Println("Hmm...schedule didn't complete")
	// }
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
	rs.cloneSchedule(f, st)    // schedule pending clones
	rs.cloneAsync(f, st, ti)   // clone missing repos (async)
	rs.cloneSummary(f, st, ti) // print summary
	rs.infoPrint(f, st)        // print startup
	rs.infoAsync(f, ti)        // get info for all repos (async)
	rs.infoSummary(f, st, ti)  // print summary
}

// VerifyChanges ...
func (rs Repos) VerifyChanges(f flags.Flags, st *stat.Stat, ti *timer.Timer) {
	rs.promptUser(f)
	rs.submitChanges(f, st, ti)
}
