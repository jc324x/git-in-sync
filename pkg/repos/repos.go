// Package repos implements access to Git repositories for git-in-sync.
package repos

import (
	"bytes"
	"fmt"
	"log"
	"sort"
	"strconv"
	"sync"
	"time"

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

// Workspaces returns the names of all Workspaces,
// reduced to single entries.
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

// Init ...
func Init(c conf.Config, f flags.Flags, t *timer.Timer) (rs Repos) {

	brf.Printv(f, "%v parsing workspaces|repos", e.Get("Pager"))

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

	brf.Printv(f, "%v [%v|%v] workspaces|repos {%v / %v}", e.Get("FaxMachine"), len(ws), len(rs), t.Split(), t.Time())

	return rs
}

func initPendingRepos(rs Repos) (prs Repos) {
	for _, r := range rs {
		if r.Category == "Pending" {
			prs = append(prs, r)
		}
	}
	return prs
}

func initScheludedRepos(rs Repos) (srs Repos) {
	for _, r := range rs {
		if r.Category == "Scheduled" {
			srs = append(srs, r)
		}
	}
	return srs
}

func initSkippedRepos(rs Repos) (skrs Repos) {
	for _, r := range rs {
		if r.Category == "Skipped" {
			skrs = append(skrs, r)
		}
	}
	return skrs
}

// sort A-Z by r.Name

// NameSort ...
func (rs Repos) NameSort() {
	sort.SliceStable(rs, func(i, j int) bool { return rs[i].Name < rs[j].Name })
}

// sort A-Z by r.WorkspacePath, then r.Name

// WorkspacePathSort  ...
func (rs Repos) WorkspacePathSort() {
	sort.SliceStable(rs, func(i, j int) bool { return rs[i].Name < rs[j].Name })
	sort.SliceStable(rs, func(i, j int) bool { return rs[i].WorkspacePath < rs[j].WorkspacePath })
}

//func (rs Repos) submitChanges(e Emoji, f Flags) {
//	srs := initScheludedRepos(rs)
//	skrs := initSkippedRepos(rs)

//	// nothing to see here, return early
//	if len(srs) == 0 && len(skrs) == 0 {
//		return
//	}

//	var wg sync.WaitGroup
//	for i := range srs {
//		wg.Add(1)
//		go func(r *Repo) {
//			defer wg.Done()
//			switch r.GitAction {
//			case "pull":
//				r.gitPull(e, f)
//			case "push":
//				r.gitPush(e, f)
//			case "add-commit-push":
//				r.gitAdd(e, f)
//				r.gitCommit(e, f)
//				r.gitPush(e, f)
//			case "stash-pull-pop-commit-push":
//				r.gitStash(e, f)
//				r.gitPull(e, f)
//				r.gitPop(e, f)
//				r.gitCommit(e, f)
//				r.gitPush(e, f)
//			}
//			r.gitRemoteUpdate(e, f)
//			r.gitStatusPorcelain(e, f)

//		}(srs[i])
//	}
//	wg.Wait()

//	var vc []string // verified complete repos

//	for _, r := range srs {
//		if r.Category == "Complete" {
//			vc = append(vc, r.Name)
//		}
//	}

//	//
//	switch {
//	case len(srs) == len(vc) && len(skrs) == 0:
//		fmt.Println("all good. nothing skipped, everything completed")
//	// case len(srs) == len(vc) && len(skrs) >= 1:
//	// 	fmt.Println("all pending actions complete - did skip this though (as planned)")
//	case len(srs) != len(vc) && len(skrs) >= 1:
//		fmt.Println("all changes not submitted correctly, also skipped")
//	}

//	// if len(srs) == len(vc) {
//	// 	fmt.Println("All changes submitted for pending repos")
//	// } else {
//	// 	fmt.Println("Hmm...schedule didn't complete")
//	// }
//}

// func (rs Repos) verifyChanges(e Emoji, f Flags) {

// 	prs := initPendingRepos(rs)

// 	if len(prs) >= 1 {
// 		for _, r := range prs {

// 			var b bytes.Buffer

// 			switch r.Status {
// 			case "Ahead":
// 				b.WriteString(e.Bunny)
// 				b.WriteString(" ")
// 				b.WriteString(r.Name)
// 				b.WriteString(" is ahead of ")
// 				b.WriteString(r.UpstreamBranch)
// 			case "Behind":
// 				b.WriteString(e.Turtle)
// 				b.WriteString(" ")
// 				b.WriteString(r.Name)
// 				b.WriteString(" is behind ")
// 				b.WriteString(r.UpstreamBranch)
// 			case "Dirty", "DirtyUntracked", "DirtyAhead", "DirtyBehind":
// 				b.WriteString(e.Pig)
// 				b.WriteString(" ")
// 				b.WriteString(r.Name)
// 				b.WriteString(" is dirty [")
// 				b.WriteString(strconv.Itoa((len(r.DiffsNameOnly))))
// 				b.WriteString("]{")
// 				b.WriteString(r.DiffsSummary)
// 				b.WriteString("}(")
// 				b.WriteString(r.ShortStatSummary)
// 				b.WriteString(")")
// 			case "Untracked", "UntrackedAhead", "UntrackedBehind":
// 				b.WriteString(e.Pig)
// 				b.WriteString(" ")
// 				b.WriteString(r.Name)
// 				b.WriteString(" is untracked [")
// 				b.WriteString(strconv.Itoa(len(r.UntrackedFiles)))
// 				b.WriteString("]{")
// 				b.WriteString(r.UntrackedSummary)
// 				b.WriteString("}")
// 			case "Up-To-Date":
// 				b.WriteString(e.Checkmark)
// 				b.WriteString(" ")
// 				b.WriteString(r.Name)
// 				b.WriteString(" is up to date with ")
// 				b.WriteString(r.UpstreamBranch)
// 			}

// 			switch r.Status {
// 			case "DirtyUntracked":
// 				b.WriteString(" and untracked [")
// 				b.WriteString(strconv.Itoa(len(r.UntrackedFiles)))
// 				b.WriteString("]{")
// 				b.WriteString(r.UntrackedSummary)
// 				b.WriteString("}")
// 			case "DirtyAhead":
// 				b.WriteString(" & ahead of ")
// 				b.WriteString(r.UpstreamBranch)
// 			case "DirtyBehind":
// 				b.WriteString(" & behind")
// 				b.WriteString(r.UpstreamBranch)
// 			case "UntrackedAhead":
// 				b.WriteString(" & is ahead of ")
// 				b.WriteString(r.UpstreamBranch)
// 			case "UntrackedBehind":
// 				b.WriteString(" & is behind ")
// 				b.WriteString(r.UpstreamBranch)
// 			}

// 			targetPrintln(f, b.String())

// 			switch r.Status {
// 			case "Ahead":
// 				fmt.Printf("%v push changes to %v? ", e.Rocket, r.Remote)
// 			case "Behind":
// 				fmt.Printf("%v pull changes from %v? ", e.Boat, r.Remote)
// 			case "Dirty":
// 				fmt.Printf("%v add all files, commit and push to %v? ", e.Clipboard, r.Remote)
// 			case "DirtyUntracked":
// 				fmt.Printf("%v add all files, commit and push to %v? ", e.Clipboard, r.Remote)
// 			case "DirtyAhead":
// 				fmt.Printf("%v add all files, commit and push to %v? ", e.Clipboard, r.Remote)
// 			case "DirtyBehind":
// 				fmt.Printf("%v stash all files, pull changes, commit and push to %v? ", e.Clipboard, r.Remote)
// 			case "Untracked":
// 				fmt.Printf("%v add all files, commit and push to %v? ", e.Clipboard, r.Remote)
// 			case "UntrackedAhead":
// 				fmt.Printf("%v add all files, commit and push to %v? ", e.Clipboard, r.Remote)
// 			case "UntrackedBehind":
// 				fmt.Printf("%v stash all files, pull changes, commit and push to %v? ", e.Clipboard, r.Remote)
// 			}

// 			// prompt for approval
// 			r.checkConfirmed()

// 			// prompt for commit message
// 			if r.Category != "Skipped" && strings.Contains(r.GitAction, "commit") {
// 				fmt.Printf("%v commit message: ", e.Memo)
// 				r.checkCommitMessage()
// 			}
// 		}

// 		// t.MarkMoment("verify-changes")

// 		// FLAG:
// 		// check again see how many pending remain, should be zero...
// 		// going to push pause for now
// 		// I need to know count of pending/scheduled prior to the start
// 		// to see what the difference is since then.
// 		// things can be autoscheduled, need to account for those

// 		// var sr []string // scheduled repos
// 		// for _, r := range rs {
// 		// 	if r.Category == "Scheduled " {
// 		// 		sr = append(sr, r.Name)
// 		// 	}
// 		// }

// 		// var b bytes.Buffer
// 		// tr := time.Millisecond // truncate

// 		// debug
// 		// for _, r := range prs {
// 		// 	fmt.Println(r.Name)
// 		// }

// 		// switch {
// 		// case len(prs) >= 1 && len(sr) >= 1:
// 		// 	b.WriteString(e.Hourglass)
// 		// 	b.WriteString(" [")
// 		// 	b.WriteString(strconv.Itoa(len(prs)))
// 		// case len(prs) >= 1 && len(sr) == 0:
// 		// 	b.WriteString(e.Warning)
// 		// 	b.WriteString(" [")
// 		// 	b.WriteString(strconv.Itoa(len(fcp)))
// 		// }

// 		// if len(prs) >= 1 && len(sr) >= 1 {
// 		// 	b.WriteString(e.Hourglass)
// 		// 	b.WriteString(" [")
// 		// 	b.WriteString(strconv.Itoa(len(prs)))
// 		// } else {
// 		// fmt.Println()
// 		// b.WriteString(e.Warning)
// 		// b.WriteString(" [")
// 		// b.WriteString(strconv.Itoa(len(fcp)))
// 		// }

// 		// b.WriteString("/")
// 		// b.WriteString(strconv.Itoa(len(prs)))
// 		// b.WriteString("] scheduled {")
// 		// b.WriteString(t.GetSplit().Truncate(tr).String())
// 		// b.WriteString(" / ")
// 		// b.WriteString(t.GetTime().Truncate(tr).String())
// 		// b.WriteString("}")

// 		// targetPrintln(f, b.String())
// 	}

// }

// VerifyWorkspaces ...
func (rs Repos) VerifyWorkspaces(f flags.Flags, ru *run.Run, t *timer.Timer) {

	// sort Repos A-Z by *r.WorkspacePath
	rs.WorkspacePathSort()

	// []string of *r.Workspace
	ws := rs.Workspaces()

	brf.Printv(f, "%v  verifying workspaces [%v](%v)", e.Get("FileCabinet"), len(ws), brf.Summary(ws, 25))

	for _, r := range rs {
		r.VerifyWorkspace(f, ru)
	}

	// summary
	var b bytes.Buffer

	if len(ws) == ru.VWC {
		b.WriteString(e.Get("Briefcase"))
	} else {
		b.WriteString(e.Get("Slash"))
	}

	b.WriteString(fmt.Sprintf(" [%v/%v] divs verified", ru.VWC, len(ws)))

	if ru.CWC >= 1 {
		b.WriteString(fmt.Sprintf(", created [%v]", strconv.Itoa(ru.CWC)))
	}

	b.WriteString(fmt.Sprintf(" {%v/%v}", t.Split().String(), t.Time().String()))

	brf.Printv(f, b.String())
}

// VerifyCloned ...
// func (rs Repos) VerifyCloned(f flags.Flags, t *timer.Timer) {

// 	var pc []string // pending clone

// 	for _, r := range rs {
// 		r.gitCheckPending()

// 		if r.PendingClone == true {
// 			pc = append(pc, r.Name)
// 		}
// 	}

// 	// return if there are no pending repos

// 	if len(pc) <= 0 {
// 		return
// 	}

// 	// if there are pending repos
// 	brf.Printv(f, "%v cloning [%v]", e.Get("Sheep"), len(pc))

// 	// verify each repo (async)
// 	var wg sync.WaitGroup
// 	for i := range rs {
// 		wg.Add(1)
// 		go func(r *Repo) {
// 			defer wg.Done()
// 			r.gitClone(f)
// 		}(rs[i])
// 	}
// 	wg.Wait()

// 	var cr []string // cloned repos

// 	for _, r := range rs {
// 		if r.Cloned == true {
// 			cr = append(cr, r.Name)
// 		}
// 	}

// 	// timer
// 	t.Mark("verify-repos")

// 	// summary
// 	var b bytes.Buffer

// 	b.WriteString(e.Get("Truck"))
// 	b.WriteString(" [")
// 	b.WriteString(strconv.Itoa(len(cr)))
// 	b.WriteString("/")
// 	b.WriteString(strconv.Itoa(len(pc)))
// 	b.WriteString("] cloned")

// 	tr := time.Millisecond // truncate
// 	b.WriteString(" {")
// 	b.WriteString(t.Split().Truncate(tr).String())
// 	b.WriteString(" / ")
// 	b.WriteString(t.Time().Truncate(tr).String())
// 	b.WriteString("}")

// 	brf.Printv(f, b.String())
// }

// VerifyRepos ...
func (rs Repos) VerifyRepos(f flags.Flags, t *timer.Timer) {
	var rn []string // repo names

	for _, r := range rs {
		rn = append(rn, r.Name)
	}

	rns := brf.Summary(rn, 25)

	// print
	brf.Printv(f, "%v  verifying repos [%v](%v)", e.Get("Satellite"), len(rs), rns)

	// verify each repo (async)
	var wg sync.WaitGroup
	for i := range rs {
		wg.Add(1)
		go func(r *repo.Repo) {
			defer wg.Done()
			// r.gitConfigOriginURL(e, f)
			// r.gitRemoteUpdate(e, f)
			// r.gitAbbrevRef(e, f)
			// r.gitLocalSHA(e, f)
			// r.gitUpstreamSHA(e, f)
			// r.gitMergeBaseSHA(e, f)
			// r.gitRevParseUpstream(e, f)
			// r.gitDiffsNameOnly(e, f)
			// r.gitShortstat(e, f)
			// r.gitUntracked(e, f)
			// r.setStatus(e, f)
		}(rs[i])
	}
	wg.Wait()

	// track Complete, Pending, Skipped and Scheduled
	var cr []string  // complete repos
	var pr []string  // pending repos
	var sk []string  // skipped repos
	var sch []string // scheduled repos

	for _, r := range rs {
		if r.Category == "Complete" {
			cr = append(cr, r.Name)
		}

		if r.Category == "Pending" {
			pr = append(pr, r.Name)
		}

		if r.Category == "Skipped" {
			sk = append(sk, r.Name)
		}

		if r.Category == "Scheduled" {
			sch = append(sch, r.Name)
		}
	}

	// timer
	t.Mark("verify-repos")

	var b bytes.Buffer

	if len(cr) == len(rs) {
		b.WriteString(e.Get("Checkmark"))
	} else {
		b.WriteString(e.Get("Traffic"))
	}

	b.WriteString(" [")
	b.WriteString(strconv.Itoa(len(cr)))
	b.WriteString("/")
	b.WriteString(strconv.Itoa(len(rs)))
	b.WriteString("] complete {")

	tr := time.Millisecond // truncate
	b.WriteString(t.Split().Truncate(tr).String())
	b.WriteString(" / ")
	b.WriteString(t.Time().Truncate(tr).String())
	b.WriteString("}")

	brf.Printv(f, b.String())

	// scheduled repo info

	if len(sch) >= 1 {
		b.Reset()
		schs := brf.Summary(sch, 15) // scheduled repo summary
		b.WriteString(e.Get("TimerClock"))
		b.WriteString("  [")
		b.WriteString(strconv.Itoa(len(sch)))

		// if flags.loginMode(f) {
		// 	b.WriteString("] pull scheduled (")

		// } else if logoutMode(f) {
		// 	b.WriteString("] push scheduled (")
		// }

		b.WriteString(schs)
		b.WriteString(")")
		brf.Printv(f, b.String())
	}

	// skipped repo info
	if len(sk) >= 1 {
		b.Reset()
		sks := brf.Summary(sk, 15) // skipped repo summary
		b.WriteString(e.Get("Slash"))
		b.WriteString(" [")
		b.WriteString(strconv.Itoa(len(sk)))
		b.WriteString("] skipped (")
		b.WriteString(sks)
		b.WriteString(")")
		brf.Printv(f, b.String())
	}

	// pending repo info
	if len(pr) >= 1 {
		b.Reset()
		prs := brf.Summary(pr, 15) // pending repo summary
		b.WriteString(e.Get("Warning"))
		b.WriteString(" [")
		b.WriteString(strconv.Itoa(len(pr)))
		b.WriteString("] pending (")
		b.WriteString(prs)
		b.WriteString(")")
		brf.Printv(f, b.String())
	}
}
