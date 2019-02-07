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
	ru.TWS = rs.Workspaces()

	ru.Reduce()

	// "verifying workspaces ..."
	brf.Printv(f, "%v  verifying workspaces [%v](%v)", e.Get("FileCabinet"), ru.TWC, brf.Summary(ru.TWS, 25))

	for _, r := range rs {
		r.VerifyWorkspace(f, ru)
	}

	ru.VWSummary(f, t)
}

// AsyncClone ...
func (rs Repos) AsyncClone(f flags.Flags, ru *run.Run) {
	var wg sync.WaitGroup
	for i := range rs {
		wg.Add(1)
		go func(r *repo.Repo) {
			defer wg.Done()
			r.GitClone(f, ru)
		}(rs[i])
	}
	wg.Wait()
}

// AsyncCheck ...
func (rs Repos) AsyncCheck() {
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

// Pulled from master branch
// func (r *Repo) verify(e Emoji, f Flags, w *Workspace) {

// 	if noDiv(r) {
// 		targetPrint(f, "%v %v (%v/%v)", e.Slash, r.Name, r.Div.Name, r.Remote)
// 		return
// 	}

// 	if canClone(f, r) {
// 		r.gitClone(e, f, w)
// 	}

// 	if isMissing(r) {
// 		targetPrint(f, "%v %v (%v/%v)", e.Slash, r.Name, r.Div.Name, r.Remote)
// 		return
// 	}

// 	r.gitConfigOriginURL()
// 	r.gitRemoteUpdate()
// 	r.gitStatusPorcelain()
// 	r.gitLocalSHA()
// 	r.gitLocalBranch()
// 	r.gitMergeBaseSHA()
// 	r.gitUpstreamSHA()
// 	r.gitRevParseUpstream()
// 	r.gitDiffsNameOnly()
// 	r.getDiffSummary()
// 	r.gitShortstat()
// 	r.getShortInts()
// 	r.gitUntracked()
// 	r.getUntrackedSummary()
// 	r.getUpstreamStatus()
// 	r.getPhase()

// 	if isUpToDate(r) {
// 		w.VerifiedRepos = append(w.VerifiedRepos, r)
// 		targetPrint(f, "%v %v (%v/%v)", e.Checkmark, r.Name, r.Div.Name, r.Remote)
// 	} else {
// 		targetPrint(f, "%v %v (%v/%v)", e.Warning, r.Name, r.Div.Name, r.Remote)
// 	}
// }

// VerifyRepos ...
func (rs Repos) VerifyRepos(f flags.Flags, ru *run.Run, t *timer.Timer) {

	for _, r := range rs {
		r.VerifyRepo(f, ru)
	}

	ru.Reduce()

	// async clone
	if ru.PCC > 1 {
		brf.Printv(f, "%v cloning [%v]", e.Get("Sheep"), ru.PCC) // "cloning..."
		rs.AsyncClone(f, ru)                                     // async clone
		ru.Reduce()                                              // reduce run
		t.Mark("async-clone")
		ru.VCSummary(f, t)
	}

	rs.AsyncCheck()

	for _, r := range rs {
		fmt.Println(r.OriginURL)
	}

}
