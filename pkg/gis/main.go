package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	// "os/exec"
	"os/user"
	// "path" -> Clean
	"sort"
	"strings"
	"sync"

	"github.com/jychri/git-in-sync/pkg/conf"
	"github.com/jychri/git-in-sync/pkg/emoji"
	"github.com/jychri/git-in-sync/pkg/flags"
	"github.com/jychri/git-in-sync/pkg/repo"
	"github.com/jychri/git-in-sync/pkg/timer"
)

// initRepo returns a *Repo with initial values set.

// --> Repos: Collection of Repos

type Repos []*repo.Repo

func initRepos(c conf.Config) (rs Repos) {

	// print
	// targetPrintln(f, "%v parsing divs|repos", e.Pager)
	// emoji.Eprint("%v parsing divs|repos", "Boat")

	// initialize Repos from Config
	for _, bl := range c.Bundles {
		for _, z := range bl.Zones {
			for _, rn := range z.Repos {
				r := repo.Init(z.Workspace, z.User, z.Remote, bl.Path, rn)
				rs = append(rs, r)
			}
		}
	}

	// timer
	// t.MarkMoment("init-repos")

	// sort
	rs.sortByPath()

	// get all divs, remove duplicates
	var dvs []string // divs

	for _, r := range rs {
		dvs = append(dvs, r.DivPath)
	}

	dvs = removeDuplicates(dvs)

	// print
	// targetPrintln(f, "%v [%v|%v] divs|repos {%v / %v}", e.FaxMachine, len(dvs), len(rs), t.GetSplit(), t.GetTime())

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
func (rs Repos) sortByName() {
	sort.SliceStable(rs, func(i, j int) bool { return rs[i].Name < rs[j].Name })
}

// sort A-Z by r.DivPath, then r.Name
func (rs Repos) sortByPath() {
	sort.SliceStable(rs, func(i, j int) bool { return rs[i].Name < rs[j].Name })
	sort.SliceStable(rs, func(i, j int) bool { return rs[i].DivPath < rs[j].DivPath })
}

// Utility functions.

func tprintln(f flags.Flags, s string, z ...interface{}) {

	if f.Check("oneline") {
		return
	}

	fmt.Println(fmt.Sprintf(s, z...))
}

func noPermission(info os.FileInfo) bool {

	if info == nil {
		return false
	}

	if len(info.Mode().String()) <= 4 {
		return true
	}

	if s := info.Mode().String()[1:4]; s != "rwx" {
		return true
	}

	return false
}

func isDirectory(info os.FileInfo) bool {
	if info == nil {
		return false
	}

	if info.IsDir() {
		return true
	} else {
		return false
	}
}

func isEmpty(p string) bool {
	f, err := os.Open(p)

	if err != nil {
		return false
	}

	_, err = f.Readdir(1)

	if err == io.EOF {
		return true
	}

	return false
}

func notEmpty(p string) bool {
	f, err := os.Open(p)

	if err != nil {
		return false
	}

	_, err = f.Readdir(1)

	if err == io.EOF {
		return false
	}

	return true
}

func isFile(info os.FileInfo) bool {
	if info == nil {
		return false
	}

	if info.IsDir() {
		return false
	} else {
		return true
	}
}

func validatePath(p string) string {
	if t := strings.TrimPrefix(p, "~/"); t != p {
		u, err := user.Current()

		if err != nil {
			log.Fatalf("Unable to identify the current user")
		}

		t := strings.Join([]string{u.HomeDir, "/", t}, "")
		return strings.TrimSuffix(t, "/")
	}
	return strings.TrimSuffix(p, "/")
}

func lastPathSelection(p string) string {
	if strings.Contains(p, "/") == true {
		sp := strings.SplitAfter(p, "/") // split path
		lp := sp[len(sp)-1]              // last path
		return lp
	} else {
		return p
	}
}

func removeDuplicates(ssl []string) (sl []string) {

	smap := make(map[string]bool)

	for i := range ssl {
		if smap[ssl[i]] == true {
		} else {
			smap[ssl[i]] = true
			sl = append(sl, ssl[i])
		}
	}

	return sl
}

func firstLine(s string) string {
	lines := strings.Split(strings.TrimSuffix(s, "\n"), "\n")

	if len(lines) >= 1 {
		return lines[0]
	} else {
		return ""
	}
}

func sliceSummary(sl []string, l int) string {
	if len(sl) == 0 {
		return ""
	}

	var csl []string // check slice
	var b bytes.Buffer

	for _, s := range sl {
		lc := len(strings.Join(csl, ", ")) // (l)ength(c)heck
		switch {
		case lc <= l-10 && len(s) <= 20: //
			csl = append(csl, s)
		case lc <= l && len(s) <= 12:
			csl = append(csl, s)
		}
	}

	b.WriteString(strings.Join(csl, ", "))

	if len(sl) != len(csl) {
		b.WriteString("...")
	}

	return b.String()
}

// --> main fns

func initRun() (f flags.Flags, rs Repos, t *timer.Timer) {

	// clear the screen
	emoji.ClearScreen()

	// initialize Timer and Flags
	t = timer.Init()
	f = flags.Init()

	// targetPrint prints a message with or without an emoji if f.Emoji is true or false.
	tprintln(f, "%v start", emoji.Eprintln("Clapper"))

	// print flag init
	// if ft, err := t.GetMoment("init-flags"); err == nil {
	// 	targetPrintln(f, "%v parsing flags", e.FlagInHole)
	// 	targetPrintln(f, "%v [%v] flags (%v) {%v / %v}", e.Flag, f.Count, f.Summary, ft.Split, ft.Start)
	// }

	// print emoji init
	// if et, err := t.GetMoment("init-emoji"); err == nil {
	// 	targetPrintln(f, "%v initializing emoji", e.CrystalBall)
	// 	targetPrintln(f, "%v [%v] emoji {%v / %v}", e.DirectHit, e.Count, et.Split, et.Start)
	// }

	// read ~/.gisrc.json, initialize Config
	// c := initConfig(e, f, t)

	// timer
	// t.MarkMoment("init-config")

	// print
	// targetPrintln(f, "%v read %v {%v / %v}", e.Book, g, t.GetSplit(), t.GetTime())

	// initialize Repos
	// rs = initRepos(c, e, f, t)

	return f, rs, t
}

// func (rs Repos) verifyDivs(f Flags) {

// 	// sort
// 	rs.sortByPath()

// 	// get all divs, remove duplicates
// 	var dvs []string  // divs
// 	var zdvs []string // zone divisions (go, main, google-apps-script etc)

// 	for _, r := range rs {
// 		dvs = append(dvs, r.DivPath)
// 		zdvs = append(zdvs, r.Division)
// 	}

// 	dvs = removeDuplicates(dvs)
// 	zdvs = removeDuplicates(zdvs)

// 	zds := sliceSummary(zdvs, 25) // zone division summary

// 	// print
// 	// targetPrintln(f, "%v  verifying divs [%v](%v)", e.FileCabinet, len(dvs), zds)

// 	// track created, verified and missing divs
// 	var cd []string // created divs
// 	var vd []string // verified divs
// 	var id []string // inaccessible divs // --> FLAG: change to unverified?

// 	for _, r := range rs {

// 		_, err := os.Stat(r.DivPath)

// 		// create div if missing and active run
// 		if os.IsNotExist(err) {
// 			// targetPrintln(f, "%v creating %v", e.Folder, r.DivPath)
// 			os.MkdirAll(r.DivPath, 0777)
// 			cd = append(cd, r.DivPath)
// 		}

// 		// check div status
// 		info, err := os.Stat(r.DivPath)

// 		switch {
// 		case noPermission(info):
// 			// r.markError(e, f, "fatal: No permsission", "verify-divs")
// 			id = append(id, r.DivPath)
// 		case !info.IsDir():
// 			// r.markError(e, f, "fatal: File occupying path", "verify-divs")
// 			id = append(id, r.DivPath)
// 		case os.IsNotExist(err):
// 			// r.markError(e, f, "fatal: No directory", "verify-divs")
// 			id = append(id, r.DivPath)
// 		case err != nil:
// 			// r.markError(e, f, "fatal: No directory", "verify-divs")
// 			id = append(id, r.DivPath)
// 		default:
// 			r.Verified = true
// 			vd = append(vd, r.DivPath)
// 		}
// 	}

// 	// timer
// 	// t.MarkMoment("verify-divs")

// 	// remove duplicates from slices
// 	vd = removeDuplicates(vd)
// 	id = removeDuplicates(id)

// 	// summary
// 	var b bytes.Buffer

// 	if len(dvs) == len(vd) {
// 		// b.WriteString(e.Briefcase)
// 	} else {
// 		// b.WriteString(e.Slash)
// 	}

// 	b.WriteString(" [")
// 	b.WriteString(strconv.Itoa(len(vd)))
// 	b.WriteString("/")
// 	b.WriteString(strconv.Itoa(len(dvs)))
// 	b.WriteString("] divs verified")

// 	if len(cd) >= 1 {
// 		b.WriteString(", created [")
// 		b.WriteString(strconv.Itoa(len(cd)))
// 		b.WriteString("]")
// 	}

// 	// b.WriteString(" {")
// 	// b.WriteString(t.GetSplit().String())
// 	// b.WriteString(" / ")
// 	// b.WriteString(t.GetTime().String())
// 	// b.WriteString("}")

// 	// targetPrintln(f, b.String())
// }

func (rs Repos) verifyCloned() {
	var pc []string // pending clone

	for _, r := range rs {
		// r.gitCheckPending(e, f)

		if r.PendingClone == true {
			pc = append(pc, r.Name)
		}
	}

	// return if there are no pending repos

	if len(pc) <= 0 {
		return
	}

	// if there are pending repos
	// targetPrintln(f, "%v cloning [%v]", e.Sheep, len(pc))

	// verify each repo (async)
	var wg sync.WaitGroup
	for i := range rs {
		wg.Add(1)
		go func(r *repo.Repo) {
			defer wg.Done()
			// r.gitClone(e, f)
		}(rs[i])
	}
	wg.Wait()

	var cr []string // cloned repos

	for _, r := range rs {
		if r.Cloned == true {
			cr = append(cr, r.Name)
		}
	}

	// timer
	// t.MarkMoment("verify-repos")

	// summary
	// var b bytes.Buffer

	// b.WriteString(e.Truck)
	// b.WriteString(" [")
	// b.WriteString(strconv.Itoa(len(cr)))
	// b.WriteString("/")
	// b.WriteString(strconv.Itoa(len(pc)))
	// b.WriteString("] cloned")

	// tr := time.Millisecond // truncate
	// b.WriteString(" {")
	// b.WriteString(t.GetSplit().Truncate(tr).String())
	// b.WriteString(" / ")
	// b.WriteString(t.GetTime().Truncate(tr).String())
	// b.WriteString("}")

	// targetPrintln(f, b.String())
}

// // func (rs Repos) verifyRepos(e Emoji, f Flags) {
// 	var rn []string // repo names

// 	for _, r := range rs {
// 		rn = append(rn, r.Name)
// 	}

// 	rns := sliceSummary(rn, 25)

// 	// print
// 	targetPrintln(f, "%v  verifying repos [%v](%v)", e.Satellite, len(rs), rns)

// 	// verify each repo (async)
// 	var wg sync.WaitGroup
// 	for i := range rs {
// 		wg.Add(1)
// 		go func(r *Repo) {
// 			defer wg.Done()
// 			r.gitConfigOriginURL(e, f)
// 			r.gitRemoteUpdate(e, f)
// 			r.gitAbbrevRef(e, f)
// 			r.gitLocalSHA(e, f)
// 			r.gitUpstreamSHA(e, f)
// 			r.gitMergeBaseSHA(e, f)
// 			r.gitRevParseUpstream(e, f)
// 			r.gitDiffsNameOnly(e, f)
// 			r.gitShortstat(e, f)
// 			r.gitUntracked(e, f)
// 			r.setStatus(e, f)
// 		}(rs[i])
// 	}
// 	wg.Wait()

// 	// track Complete, Pending, Skipped and Scheduled
// 	var cr []string  // complete repos
// 	var pr []string  // pending repos
// 	var sk []string  // skipped repos
// 	var sch []string // scheduled repos

// 	for _, r := range rs {
// 		if r.Category == "Complete" {
// 			cr = append(cr, r.Name)
// 		}

// 		if r.Category == "Pending" {
// 			pr = append(pr, r.Name)
// 		}

// 		if r.Category == "Skipped" {
// 			sk = append(sk, r.Name)
// 		}

// 		if r.Category == "Scheduled" {
// 			sch = append(sch, r.Name)
// 		}
// 	}

// 	// timer
// 	// t.MarkMoment("verify-repos")

// 	// var b bytes.Buffer

// 	// if len(cr) == len(rs) {
// 	// 	b.WriteString(e.Checkmark)
// 	// } else {
// 	// 	b.WriteString(e.Traffic)
// 	// }

// 	// b.WriteString(" [")
// 	// b.WriteString(strconv.Itoa(len(cr)))
// 	// b.WriteString("/")
// 	// b.WriteString(strconv.Itoa(len(rs)))
// 	// b.WriteString("] complete {")

// 	// tr := time.Millisecond // truncate
// 	// b.WriteString(t.GetSplit().Truncate(tr).String())
// 	// b.WriteString(" / ")
// 	// b.WriteString(t.GetTime().Truncate(tr).String())
// 	// b.WriteString("}")

// 	// targetPrintln(f, b.String())

// 	// scheduled repo info

// 	// if len(sch) >= 1 {
// 	// 	b.Reset()
// 	// 	schs := sliceSummary(sch, 15) // scheduled repo summary
// 	// 	b.WriteString(e.TimerClock)
// 	// 	b.WriteString("  [")
// 	// 	b.WriteString(strconv.Itoa(len(sch)))

// 	// 	if loginMode(f) {
// 	// 		b.WriteString("] pull scheduled (")

// 	// 	} else if logoutMode(f) {
// 	// 		b.WriteString("] push scheduled (")
// 	// 	}

// 	// 	b.WriteString(schs)
// 	// 	b.WriteString(")")
// 	// 	targetPrintln(f, b.String())
// 	// }

// 	// skipped repo info
// 	// if len(sk) >= 1 {
// 	// 	b.Reset()
// 	// 	sks := sliceSummary(sk, 15) // skipped repo summary
// 	// 	b.WriteString(e.Slash)
// 	// 	b.WriteString(" [")
// 	// 	b.WriteString(strconv.Itoa(len(sk)))
// 	// 	b.WriteString("] skipped (")
// 	// 	b.WriteString(sks)
// 	// 	b.WriteString(")")
// 	// 	targetPrintln(f, b.String())
// 	// }

// 	// pending repo info
// 	// if len(pr) >= 1 {
// 	// 	b.Reset()
// 	// 	prs := sliceSummary(pr, 15) // pending repo summary
// 	// 	b.WriteString(e.Warning)
// 	// 	b.WriteString(" [")
// 	// 	b.WriteString(strconv.Itoa(len(pr)))
// 	// 	b.WriteString("] pending (")
// 	// 	b.WriteString(prs)
// 	// 	b.WriteString(")")
// 	// 	targetPrintln(f, b.String())
// 	// }

// }

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

// FLAG: need to fix up messaging here
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

// debug spits out error info
func (rs Repos) debug() {
	for _, r := range rs {
		if r.ErrorShort != "" {
			fmt.Printf("%v|%v (%v)\n", r.Name, r.ErrorName, r.ErrorFirst)
			fmt.Printf("%v\n", r.ErrorShort)
			fmt.Printf("clean: %v, untracked: %v, status: %v\n", r.Clean, r.Untracked, r.Status)
		}
	}
}

func main() {
	// f, rs, t := initRun()
	initRun()
	// fmt.Println(f)
	// rs.verifyDivs(e, f)
	// rs.verifyCloned(e, f)
	// rs.verifyRepos(e, f)
	// rs.verifyChanges(e, f)
	// rs.submitChanges(e, f)
	// rs.debug()
}
