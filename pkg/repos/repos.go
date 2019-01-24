package repos

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/jychri/git-in-sync/pkg/brf"
	"github.com/jychri/git-in-sync/pkg/conf"
	"github.com/jychri/git-in-sync/pkg/emoji"
	"github.com/jychri/git-in-sync/pkg/fchk"
	"github.com/jychri/git-in-sync/pkg/flags"
)

type Repo struct {

	// initRun -> initRepos -> initRepo
	BundlePath string // "~/dev"
	Division   string // "main" or "go-lang"
	User       string // "jychri"
	Remote     string // "github" or "gitlab"
	Name       string // "git-in-sync"
	DivPath    string // "/Users/jychri/dev/go-lang/"
	RepoPath   string // "/Users/jychri/dev/go-lang/git-in-sync"
	GitPath    string // "/Users/jychri/dev/go-lang/git-in-sync/.git"
	GitDir     string // "--git-dir=/Users/jychri/dev/go-lang/git-in-sync/.git"
	WorkTree   string // "--work-tree=/Users/jychri/dev/go-lang/git-in-sync"
	URL        string // "https://github.com/jychri/git-in-sync"

	// rs.verifyRepos
	PendingClone bool // true if RepoPath or GitPath are empty

	// rs.verifyDivs, rs.verifyRepos
	Verified     bool   // true if Repo continues to pass verification
	ErrorMessage string // the last error message
	ErrorName    string // name of the last error
	ErrorFirst   string // first line of the last error message
	ErrorShort   string // message in matched short form

	// rs.verifyRepos -> gitVerify -> gitClone
	Cloned bool // true if Repo was cloned

	// rs.verifyRepos -> gitConfigOriginURL
	OriginURL string // "https://github.com/jychri/git-in-sync"

	// rs.verifyRepos -> gitAbbrevRef
	LocalBranch string // `git rev-parse --abbrev-ref HEAD`, "master"

	// rs.verifyRepos -> gitLocalSHA
	LocalSHA string // `git rev-parse @`, "l00000ngSHA1slong324"

	// rs.verifyRepos -> gitUpstreamSHA
	UpstreamSHA string // `git rev-parse @{u}`, "l00000ngSHA1slong324"

	// rs.verifyRepos -> gitMergeBaseSHA
	MergeSHA string // `git merge-base @ @{u}`, "l00000ngSHA1slong324"

	// rs.verifyRepos -> gitRevParseUpstream
	UpstreamBranch string // `git rev-parse --abbrev-ref --symbolic-full-name @{u}`, "..."

	// rs.verifyRepos -> gitDiffsNameOnly
	DiffsNameOnly []string // `git diff --name-only @{u}`, [a, b, c, d, e]
	DiffsSummary  string   // "a, b, c..."

	// rs.verifyRepos -> gitShortstat
	ShortStat        string // `git diff --shortstat`, "x files changed, y insertions(+), z deletions(-)"
	Changed          int    // x
	Insertions       int    // y
	Deletions        int    // z
	ShortStatSummary string // "+y|-z" or "D" for Deleted if (x >= 1 && y == 0 && z == 0)
	Clean            bool   // true if Changed, Insertions and Deletions are all 0

	// rs.verifyRepos -> gitUntracked
	UntrackedFiles   []string // `git ls-files --others --exclude-standard`, [a, b, c, d, e]
	UntrackedSummary string   // "a, b, c..."
	Untracked        bool     // true if if len(r.UntrackedFiles) >= 1

	// rs.verifyRepos -> setStatus (verify grouping)
	Category   string // Complete, Pending, Skipped, Scheduled
	Status     string // better term?
	GitAction  string // "..."
	GitMessage string // "..."

	// rs.verifyChanges -> gitPorcelain
	Porcelain bool // true if `git status --porcelain` returns ""
}

func initR(zd string, zu string, zr string, bp string, rn string) *Repo {

	r := new(Repo)

	// "~/dev", (b)undle(p)ath
	r.BundlePath = bp

	// "main" or "go-lang", (z)one(d)ivision
	r.Division = zd

	// "jychri", (z)one(u)ser
	r.User = zu

	// "github" or "gitlab", (z)one(r)emote
	r.Remote = zr

	// "git-in-sync", (r)epo(n)ame
	r.Name = rn

	var b bytes.Buffer

	// "/Users/jychri/dev/go-lang/"
	// b.WriteString(validatePath(r.BundlePath))
	if r.Division != "main" {
		b.WriteString("/")
		b.WriteString(r.Division)
	}
	r.DivPath = b.String()

	// "/Users/jychri/dev/go-lang/git-in-sync/"
	b.Reset()
	b.WriteString(r.DivPath)
	b.WriteString("/")
	b.WriteString(r.Name)
	r.RepoPath = b.String()

	// "/Users/jychri/dev/go-lang/git-in-sync/.git"
	b.Reset()
	b.WriteString(r.RepoPath)
	b.WriteString("/.git")
	r.GitPath = b.String()

	// "--git-dir=/Users/jychri/dev/go-lang/git-in-sync/.git"
	b.Reset()
	b.WriteString("--git-dir=")
	b.WriteString(r.GitPath)
	r.GitDir = b.String()

	// "--work-tree=/Users/jychri/dev/go-lang/git-in-sync"
	b.Reset()
	b.WriteString("--work-tree=")
	b.WriteString(r.RepoPath)
	r.WorkTree = b.String()

	// "https://github.com/jychri/git-in-sync"
	b.Reset()
	switch r.Remote {
	case "github":
		b.WriteString("https://github.com/")
	case "gitlab":
		b.WriteString("https://gitlab.com/")
	}
	b.WriteString(r.User)
	b.WriteString("/")
	b.WriteString(r.Name)
	r.URL = b.String()

	return r
}

func notVerified(r *Repo) bool {
	if r.Verified == false {
		return true
	} else {
		return false
	}
}

func (r *Repo) markError(err string, name string) {
	r.ErrorMessage = err
	r.ErrorName = name
	// r.ErrorFirst = firstLine(err)

	if strings.Contains(r.ErrorFirst, "warning") {
		r.Verified = true
	}

	if strings.Contains(r.ErrorFirst, "fatal") {
		r.Verified = false
	}
}

func captureOut(b bytes.Buffer) string {
	return strings.TrimSuffix(b.String(), "\n")
}

func (r *Repo) gitCheckPending() {

	// return if not verified
	if notVerified(r) {
		return
	}

	// check if RepoPath and GitPath are accessible
	// rinfo, rerr := os.Stat(r.RepoPath)
	// ginfo, gerr := os.Stat(r.GitPath)

	switch {
	// case isFile(rinfo):
	// 	// r.markError(e, f, "fatal: file occupying path", "git-verify")
	// case isDirectory(rinfo) && notEmpty(r.RepoPath) && os.IsNotExist(gerr):
	// 	// r.markError(e, f, "fatal: directory occupying path", "git-verify")
	// case isDirectory(rinfo) && isEmpty(r.RepoPath):
	// 	r.PendingClone = true
	// case os.IsNotExist(rerr) && os.IsNotExist(gerr):
	// 	r.PendingClone = true
	// case isDirectory(rinfo) && isDirectory(ginfo):
	// 	r.Verified = true
	}
}

func (r *Repo) gitClone() {

	if r.PendingClone == true {
		// print
		// targetPrintln(f, "%v cloning %v {%v}", e.Box, r.Name, r.Division)

		// command
		args := []string{"clone", r.URL, r.RepoPath}
		cmd := exec.Command("git", args...)
		var out bytes.Buffer
		var err bytes.Buffer
		cmd.Stderr = &err
		cmd.Stdout = &out
		cmd.Run()

		// check error, set value(s)
		if err := err.String(); err != "" {
			// r.markError(e, f, err, "gitClone")
		}

		r.Cloned = true

	}

}

func (r *Repo) gitConfigOriginURL() {

	// return if not verified
	if notVerified(r) {
		return
	}

	// command
	args := []string{r.GitDir, "config", "--get", "remote.origin.url"}
	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()

	// trim "\n" from command output
	s := out.String()
	s = strings.TrimSuffix(s, "\n")

	// set OriginURL
	r.OriginURL = s

	// check error, set value(s)
	switch {
	case r.OriginURL == "":
		// r.markError(e, f, "fatal: 'origin' does not appear to be a git repository", "gitConfigOriginURL")
	case r.OriginURL != r.URL:
		// r.markError(e, f, "fatal: URL != OriginURL", "gitConfigOriginURL")
	}
}

func (r *Repo) gitRemoteUpdate() {

	// return if not verified
	if notVerified(r) {
		return
	}

	// command
	args := []string{r.GitDir, r.WorkTree, "fetch", "origin"}
	cmd := exec.Command("git", args...)
	var err bytes.Buffer
	cmd.Stderr = &err
	cmd.Run()

	// Warnings for redirects to "*./git" can be ignored.
	eval := err.String()
	wgit := strings.Join([]string{r.URL}, "/.git") // (w)ith .(git)

	switch {
	case strings.Contains(eval, "warning: redirecting") && strings.Contains(eval, wgit):
		// fmt.Printf("%v - redirect to .git\n", r.Name)
	case eval != "":
		// r.markError(e, f, eval, "gitRemoteUpdate")
	}
}

func (r *Repo) gitAbbrevRef() {

	// return if not verified
	if notVerified(r) {
		return
	}

	// command
	args := []string{r.GitDir, r.WorkTree, "rev-parse", "--abbrev-ref", "HEAD"}
	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	var err bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err
	cmd.Run()

	// check error, set value(s)
	if err := err.String(); err != "" {
		// r.markError(e, f, err, "gitAbbrevRef")
	} else {
		r.LocalBranch = captureOut(out)
	}
}

func (r *Repo) gitLocalSHA() {

	// return if not verified
	if notVerified(r) {
		return
	}

	// command
	args := []string{r.GitDir, r.WorkTree, "rev-parse", "@"}
	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	var err bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err
	cmd.Run()

	// check error, set value(s)
	if err := err.String(); err != "" {
		// r.markError(e, f, err, "gitLocalSHA")
	} else {
		r.LocalSHA = captureOut(out)
	}
}

func (r *Repo) gitUpstreamSHA() {

	// return if not verified
	if notVerified(r) {
		return
	}

	// command
	args := []string{r.GitDir, r.WorkTree, "rev-parse", "@{u}"}
	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	var err bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err
	cmd.Run()

	// check error, set value(s)
	if err := err.String(); err != "" {
		// r.markError(e, f, err, "gitUpstreamSHA")
	} else {
		r.UpstreamSHA = captureOut(out)
	}
}

func (r *Repo) gitMergeBaseSHA() {

	// return if not verified
	if notVerified(r) {
		return
	}

	// command
	args := []string{r.GitDir, r.WorkTree, "merge-base", "@", "@{u}"}
	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	var err bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err
	cmd.Run()

	// check error, set value(s)
	if err := err.String(); err != "" {
		// r.markError(e, f, err, "gitUpstreamSHA")
	} else {
		r.MergeSHA = captureOut(out)
	}
}

func (r *Repo) gitRevParseUpstream() {

	// return if not verified
	if notVerified(r) {
		return
	}

	// command
	args := []string{r.GitDir, r.WorkTree, "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}"}
	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	var err bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err
	cmd.Run()

	// check error, set value(s)
	if err := err.String(); err != "" {
		// r.markError(e, f, err, "gitRevParseUpstream")
	} else {
		r.UpstreamBranch = captureOut(out)
	}
}

func (r *Repo) gitDiffsNameOnly() {

	// return if not verified
	if notVerified(r) {
		return
	}

	// command
	args := []string{r.GitDir, r.WorkTree, "diff", "--name-only", "@{u}"}
	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	var err bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err
	cmd.Run()

	// check error, set value(s)
	if err := err.String(); err != "" {
		// r.markError(e, f, err, "gitDiffsNameOnly")
	}

	if str := out.String(); str != "" {
		r.DiffsNameOnly = strings.Fields(str)
		// r.DiffsSummary = sliceSummary(r.DiffsNameOnly, 12)
	} else {
		r.DiffsNameOnly = make([]string, 0)
		r.DiffsSummary = ""
	}
}

func (r *Repo) gitShortstat() {

	// return if not verified
	if notVerified(r) {
		return
	}

	// command
	args := []string{r.GitDir, r.WorkTree, "diff", "--shortstat"}
	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	var err bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err
	cmd.Run()

	// check error, set value(s)
	if err := err.String(); err != "" {
		// r.markError(e, f, err, "gitShortstat")
	} else {
		r.ShortStat = captureOut(out)
	}

	// scrape with regular expressions
	rxc := regexp.MustCompile(`(.*)? file`)
	rxs := rxc.FindStringSubmatch(r.ShortStat)
	if len(rxs) == 2 {
		s := strings.TrimPrefix(rxs[1], " ")
		if i, err := strconv.Atoi(s); err == nil {
			r.Changed = i
		}
	}

	rxi := regexp.MustCompile(`changed, (.*)? insertion`)
	rxs = rxi.FindStringSubmatch(r.ShortStat)
	if len(rxs) == 2 {
		s := rxs[1]
		if i, err := strconv.Atoi(s); err == nil {
			r.Insertions = i
		}
	}

	if r.Insertions >= 1 {
		rxd := regexp.MustCompile(`\(\+\), (.*)? deletion`)
		rxs = rxd.FindStringSubmatch(r.ShortStat)
		if len(rxs) == 2 {
			s := rxs[1]
			if i, err := strconv.Atoi(s); err == nil {
				r.Deletions = i
			}
		}
	} else {
		rxd := regexp.MustCompile(`changed, (.*)? deletion`)
		rxs = rxd.FindStringSubmatch(r.ShortStat)
		if len(rxs) == 2 {
			s := rxs[1]
			if i, err := strconv.Atoi(s); err == nil {
				r.Deletions = i
			}
		}

	}

	// set Clean and ShortStatSummary
	switch {
	case r.Changed == 0 && r.Insertions == 0 && r.Deletions == 0:
		r.Clean = true
		r.ShortStatSummary = ""
	case r.Changed >= 1 && r.Insertions == 0 && r.Deletions == 0:
		r.Clean = false
		r.ShortStatSummary = ("D")
	default:
		r.Clean = false

		var b bytes.Buffer
		b.WriteString("+")
		b.WriteString(strconv.Itoa(r.Insertions))
		b.WriteString("|-")
		b.WriteString(strconv.Itoa(r.Deletions))
		r.ShortStatSummary = b.String()
	}

	if r.Changed == 0 && r.Insertions == 0 && r.Deletions == 0 {
		r.Clean = true
	} else {
		r.Clean = false
	}

}

func (r *Repo) gitUntracked() {

	// return if not verified
	if notVerified(r) {
		return
	}

	// command
	args := []string{r.GitDir, r.WorkTree, "ls-files", "--others", "--exclude-standard"}
	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	var err bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err
	cmd.Run()

	// check error, set value(s)
	if err := err.String(); err != "" {
		// r.markError(e, f, err, "gitUntracked")
	}

	if str := out.String(); str != "" {
		ufr := strings.Fields(str) // untracked files raw
		for _, f := range ufr {
			// f = lastPathSelection(f)
			r.UntrackedFiles = append(r.UntrackedFiles, f)
			// r.UntrackedSummary = sliceSummary(r.UntrackedFiles, 12)
		}
	} else {
		r.UntrackedFiles = make([]string, 0)
	}

	if len(r.UntrackedFiles) >= 1 {
		r.Untracked = true
	}

}

func (r *Repo) setStatus() {

	switch {
	case r.LocalSHA == r.UpstreamSHA:
		r.Status = "Up-To-Date"
	case r.LocalSHA == r.MergeSHA:
		r.Status = "Behind"
	case r.UpstreamSHA == r.MergeSHA:
		r.Status = "Ahead"
	}

	switch {
	case r.Verified == false:
		r.Category = "Skipped"
		r.Status = "Error"
	case (r.Clean == true && r.Untracked == false && r.Status == "Ahead"):
		r.Category = "Pending"
		r.Status = "Ahead"
		r.GitAction = "push"
	case (r.Clean == true && r.Untracked == false && r.Status == "Behind"):
		r.Category = "Pending"
		r.Status = "Behind"
		r.GitAction = "pull"
	case (r.Clean == false && r.Untracked == false && r.Status == "Up-To-Date"):
		r.Category = "Pending"
		r.Status = "Dirty"
		r.GitAction = "add-commit-push"
	case (r.Clean == false && r.Untracked == true && r.Status == "Up-To-Date"):
		r.Category = "Pending"
		r.Status = "DirtyUntracked"
		r.GitAction = "add-commit-push"
	case (r.Clean == false && r.Untracked == false && r.Status == "Ahead"):
		r.Category = "Pending"
		r.Status = "DirtyAhead"
		r.GitAction = "add-commit-push"
	case (r.Clean == false && r.Untracked == false && r.Status == "Behind"):
		r.Category = "Pending"
		r.Status = "DirtyBehind"
		r.GitAction = "stash-pull-pop-commit-push"
	case (r.Clean == true && r.Untracked == true && r.Status == "Up-To-Date"):
		r.Category = "Pending"
		r.Status = "Untracked"
		r.GitAction = "add-commit-push"
	case (r.Clean == true && r.Untracked == true && r.Status == "Ahead"):
		r.Category = "Pending"
		r.Status = "UntrackedAhead"
		r.GitAction = "add-commit-push"
	case (r.Clean == false && r.Untracked == true && r.Status == "Behind"):
		r.Category = "Pending"
		r.Status = "UntrackedBehind"
	case (r.Clean == true && r.Untracked == false && r.Status == "Up-To-Date"):
		r.Category = "Complete"
		r.Status = "Up-To-Date"
		r.GitAction = "stash-pull-pop-commit-push"
	default:
		r.Category = "Skipped"
		// r.Status = "Unknown"
		// r.markError(e, f, "fatal: no matches found in setStatus switch", "setStatus")
	}

	if r.ErrorMessage != "" {
		err := r.ErrorMessage
		switch {
		case strings.Contains(err, "fatal: ambiguous argument 'HEAD'"):
			r.ErrorShort = "fatal: empty repository"
		case strings.Contains(err, "fatal: 'origin' does not appear to be a git repository"):
			r.ErrorShort = "fatal: 'origin' not set"
		case strings.Contains(err, "fatal: URL != OriginURL"):
			r.ErrorShort = "fatal: URL mismatch"
		case strings.Contains(err, "fatal: no matches found"):
			r.ErrorShort = "fatal: no matches found"
		}
	}

	// auto move to scheduled for matching login/logout
	// switch {
	// case loginMode(f) && r.Category == "Pending" && r.Status == "Behind":
	// 	r.Category = "Scheduled"
	// case logoutMode(f) && r.Category == "Pending" && r.Status == "Ahead":
	// 	r.Category = "Scheduled"
	// }
}

func (r *Repo) checkConfirmed() {

	// setup reader
	rdr := bufio.NewReader(os.Stdin)
	in, err := rdr.ReadString('\n')

	// return if error
	if err != nil {
		r.Category = "Skipped"
		return
	}

	// trim trailing new line
	in = strings.TrimSuffix(in, "\n")

	switch in {
	case "please", "y", "ye", "yes", "ys", "1", "ok", "push", "pull":
		r.Category = "Scheduled"
	case "you may fire when ready", "do it", "just do it", "you betcha", "sure":
		r.Category = "Scheduled"
	case "n", "no", "nah", "0", "stop", "skip", "abort", "halt", "quit":
		r.Category = "Skipped"
	default:
		r.Category = "Skipped"
	}
}

func (r *Repo) checkCommitMessage() {

	// setup reader
	rdr := bufio.NewReader(os.Stdin)
	in, err := rdr.ReadString('\n')

	// return if error
	if err != nil {
		r.Category = "Skipped"
		return
	}

	// trim trailing new line
	in = strings.TrimSuffix(in, "\n")

	switch in {
	case "n", "no", "nah", "0", "stop", "skip", "abort", "halt", "quit", "exit", "":
		r.Category = "Skipped"
		r.GitMessage = ""
	default:
		r.Category = "Scheduled"
		r.GitMessage = in
	}
}

func (r *Repo) gitAdd() {
	switch r.Status {
	case "Dirty", "DirtyUntracked", "DirtyAhead", "DirtyBehind":
		// targetPrintln(f, "%v %v adding changes [%v]{%v}(%v)", e.Outbox, r.Name, len(r.DiffsNameOnly), r.DiffsSummary, r.ShortStatSummary)
	case "Untracked", "UntrackedAhead", "UntrackedBehind":
		// targetPrintln(f, "%v %v adding new files [%v]{%v}", e.Outbox, r.Name, len(r.UntrackedFiles), r.UntrackedSummary)
	}

	// command
	args := []string{"-C", r.RepoPath, "add", "-A"}
	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	var err bytes.Buffer
	cmd.Stderr = &err
	cmd.Stdout = &out
	cmd.Run()

	// check error, set value(s)
	if err := err.String(); err != "" {
		// r.markError(e, f, err, "gitAdd")
	}

}

func (r *Repo) gitCommit() {
	switch r.Status {
	case "Dirty", "DirtyUntracked", "DirtyAhead", "DirtyBehind":
		// targetPrintln(f, "%v %v committing changes [%v]{%v}(%v)", e.Fire, r.Name, len(r.DiffsNameOnly), r.DiffsSummary, r.ShortStatSummary)
	case "Untracked", "UntrackedAhead", "UntrackedBehind":
		// targetPrintln(f, "%v %v committing new files [%v]{%v}", e.Fire, r.Name, len(r.UntrackedFiles), r.UntrackedSummary)
	}

	// command
	args := []string{"-C", r.RepoPath, "commit", "-m", r.GitMessage}
	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	var err bytes.Buffer
	cmd.Stderr = &err
	cmd.Stdout = &out
	cmd.Run()

	// check error, set value(s)
	if err := err.String(); err != "" {
		// r.markError(e, f, err, "gitCommit")
	}

}

func (r *Repo) gitStash() {
	// targetPrintln(f, "%v  %v stashing changes", e.Squirrel, r.Name)

}

func (r *Repo) gitPop() {
	// targetPrintln(f, "%v %v popping changes", e.Popcorn, r.Name)
}

func (r *Repo) gitPull() {
	// targetPrintln(f, "%v %v pulling from %v @ %v", e.Ship, r.Name, r.UpstreamBranch, r.Remote)

	// command
	args := []string{"-C", r.RepoPath, "pull"}
	cmd := exec.Command("git", args...)
	// var out bytes.Buffer
	var err bytes.Buffer
	cmd.Stderr = &err
	// cmd.Stdout = &out
	cmd.Run()

	// check error, set value(s)
	if err := err.String(); err != "" {
		// r.markError(e, f, err, "gitPull")
	}

	if r.Verified == false {
		// targetPrintln(f, "%v %v pull failed", e.Slash, r.Name)
	}
}

func (r *Repo) gitPush() {
	// targetPrintln(f, "%v %v pushing to %v @ %v", e.Rocket, r.Name, r.UpstreamBranch, r.Remote)

	// command
	args := []string{"-C", r.RepoPath, "push"}
	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	var err bytes.Buffer
	cmd.Stderr = &err
	cmd.Stdout = &out
	cmd.Run()

	// check error, set value(s)
	if err := err.String(); err != "" {
		// r.markError(e, f, err, "gitPush")
	}

	if r.Verified == false {
		// targetPrintln(f, "%v %v push failed", e.Slash, r.Name)
	}

}

func (r *Repo) gitStatusPorcelain() {

	// return if not verified
	if notVerified(r) {
		return
	}

	// command
	args := []string{r.GitDir, r.WorkTree, "status", "--porcelain"}
	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	var err bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err
	cmd.Run()

	// check error, set value(s)
	if err := err.String(); err != "" {
		// r.markError(e, f, err, "gitStatusPorcelain")
	}

	if str := out.String(); str != "" {
		r.Porcelain = false
		// targetPrintln(f, "%v commit error (%v)", e.Slash, r.ErrorFirst)
	} else {
		r.Category = "Complete"
		r.Porcelain = true
		// targetPrintln(f, "%v %v up to date!", e.Checkmark, r.Name)
	}

}

// initRepo returns a *Repo with initial values set.

// --> Repos: Collection of Repos

type Repos []*Repo

func Init(c conf.Config) (rs Repos) {

	// print
	// targetPrintln(f, "%v parsing divs|repos", e.Pager)
	// emoji.Eprint("%v parsing divs|repos", "Boat")

	// initialize Repos from Config
	for _, bl := range c.Bundles {
		for _, z := range bl.Zones {
			for _, rn := range z.Repos {
				r := initR(z.Workspace, z.User, z.Remote, bl.Path, rn)
				rs = append(rs, r)
			}
		}
	}

	// timer
	// t.MarkMoment("init-repos")

	// sort
	rs.PathSort()

	// get all divs, remove duplicates
	var dvs []string // divs

	for _, r := range rs {
		dvs = append(dvs, r.DivPath)
	}

	dvs = brf.Single(dvs)

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
func (rs Repos) NameSort() {
	sort.SliceStable(rs, func(i, j int) bool { return rs[i].Name < rs[j].Name })
}

// sort A-Z by r.DivPath, then r.Name
func (rs Repos) PathSort() {
	sort.SliceStable(rs, func(i, j int) bool { return rs[i].Name < rs[j].Name })
	sort.SliceStable(rs, func(i, j int) bool { return rs[i].DivPath < rs[j].DivPath })
}

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
		go func(r *Repo) {
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

func (rs Repos) VerifyDivs(e emoji.Emoji, f flags.Flags) {

	// sort
	rs.PathSort()

	// get all divs, remove duplicates
	var dvs []string  // divs
	var zdvs []string // zone divisions (go, main, google-apps-script etc)

	for _, r := range rs {
		dvs = append(dvs, r.DivPath)
		zdvs = append(zdvs, r.Division)
	}

	dvs = brf.Single(dvs)
	zdvs = brf.Single(zdvs)

	// zds := brf.Summary(zdvs, 25) // zone division summary

	// print
	// targetPrintln(f, "%v  verifying divs [%v](%v)", e.FileCabinet, len(dvs), zds)

	// track created, verified and missing divs
	var cd []string // created divs
	var vd []string // verified divs
	var id []string // inaccessible divs // --> FLAG: change to unverified?

	for _, r := range rs {

		_, err := os.Stat(r.DivPath)

		// create div if missing and active run
		if os.IsNotExist(err) {
			// targetPrintln(f, "%v creating %v", e.Folder, r.DivPath)
			os.MkdirAll(r.DivPath, 0777)
			cd = append(cd, r.DivPath)
		}

		// check div status
		info, err := os.Stat(r.DivPath)

		switch {
		case fchk.NoPermission(info):
			// r.markError(e, f, "fatal: No permsission", "verify-divs")
			id = append(id, r.DivPath)
		case !info.IsDir():
			// r.markError(e, f, "fatal: File occupying path", "verify-divs")
			id = append(id, r.DivPath)
		case os.IsNotExist(err):
			// r.markError(e, f, "fatal: No directory", "verify-divs")
			id = append(id, r.DivPath)
		case err != nil:
			// r.markError(e, f, "fatal: No directory", "verify-divs")
			id = append(id, r.DivPath)
		default:
			r.Verified = true
			vd = append(vd, r.DivPath)
		}
	}

	// timer
	// t.MarkMoment("verify-divs")

	// remove duplicates from slices
	vd = brf.Single(vd)
	id = brf.Single(id)

	// summary
	var b bytes.Buffer

	if len(dvs) == len(vd) {
		// b.WriteString(e.Briefcase)
	} else {
		// b.WriteString(e.Slash)
	}

	b.WriteString(" [")
	b.WriteString(strconv.Itoa(len(vd)))
	b.WriteString("/")
	b.WriteString(strconv.Itoa(len(dvs)))
	b.WriteString("] divs verified")

	if len(cd) >= 1 {
		b.WriteString(", created [")
		b.WriteString(strconv.Itoa(len(cd)))
		b.WriteString("]")
	}

	// b.WriteString(" {")
	// b.WriteString(t.GetSplit().String())
	// b.WriteString(" / ")
	// b.WriteString(t.GetTime().String())
	// b.WriteString("}")

	// targetPrintln(f, b.String())
}
