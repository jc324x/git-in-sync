// Package repo implements support for Git repositories.
package repo

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/jychri/git-in-sync/pkg/brf"
	"github.com/jychri/git-in-sync/pkg/emoji"
	"github.com/jychri/git-in-sync/pkg/fchk"
	"github.com/jychri/git-in-sync/pkg/flags"
	"github.com/jychri/git-in-sync/pkg/stat"
	"github.com/jychri/git-in-sync/pkg/tilde"
)

// private

// git runs a Git command and returns standard out
// and standard error values as strings out and em.
// em is used rather than err to indicate that the
// value is as string rather than an error.
func (r *Repo) git(args []string) (out string, em string) {

	if r.Verified == false {
		return
	}

	var outb, errb bytes.Buffer

	cmd := exec.Command("git", args...)
	cmd.Stderr = &errb
	cmd.Stdout = &outb
	cmd.Run()

	out = outb.String()
	em = errb.String()

	out = strings.TrimSuffix(out, "\n")
	em = strings.TrimSuffix(em, "\n")

	return out, em
}

// Public

// Repo models a Git repository.
type Repo struct {
	BundlePath       string   // "~/tmpgis"
	Workspace        string   // "main" or "go-lang"
	User             string   // "jychri"
	Remote           string   // "github" or "gitlab"
	Name             string   // "git-in-sync"
	WorkspacePath    string   // "/Users/jychri/tmpgis/go-lang/"
	RepoPath         string   // "/Users/jychri/tmpgis/go-lang/git-in-sync"
	GitPath          string   // "/Users/jychri/tmpgis/go-lang/git-in-sync/.git"
	GitDir           string   // "--git-dir=/Users/jychri/tmpgis/go-lang/git-in-sync/.git"
	WorkTree         string   // "--work-tree=/Users/jychri/tmpgis/go-lang/git-in-sync"
	URL              string   // "https://github.com/jychri/git-in-sync"
	PendingClone     bool     // true if RepoPath or GitPath are empty
	Verified         bool     // true if Repo continues to pass verification
	ErrorMessage     string   // the last error message
	ErrorName        string   // name of the last error
	ErrorFirst       string   // first line of the last error message
	ErrorShort       string   // message in matched short form
	Cloned           bool     // true if Repo was cloned
	OriginURL        string   // "https://github.com/jychri/git-in-sync"
	LocalBranch      string   // `git rev-parse --abbrev-ref HEAD`, "master"
	LocalSHA         string   // `git rev-parse @`, "l00000ngSHA1slong324"
	UpstreamSHA      string   // `git rev-parse @{u}`, "l00000ngSHA1slong324"
	MergeSHA         string   // `git merge-base @ @{u}`, "l00000ngSHA1slong324"
	UpstreamBranch   string   // `git rev-parse --abbrev-ref --symbolic-full-name @{u}`, "..."
	DiffsNameOnly    []string // `git diff --name-only @{u}`, [a, b, c, d, e]
	DiffsSummary     string   // "a, b, c..."
	ShortStat        string   // `git diff --shortstat`, "x files changed, y insertions(+), z deletions(-)"
	Changed          int      // x
	Insertions       int      // y
	Deletions        int      // z
	ShortStatSummary string   // "+y|-z" or "D" for Deleted if (x >= 1 && y == 0 && z == 0)
	Clean            bool     // true if Changed, Insertions and Deletions are all 0
	Untracked        bool     // true if if len(r.UntrackedFiles) >= 1
	UntrackedFiles   []string // `git ls-files --others --exclude-standard`, [a, b, c, d, e]
	UntrackedSummary string   // "a, b, c..."
	Category         string   // Complete, Pending, Skipped, Scheduled
	Status           string   // better term?
	Action           string   // "..."
	Prompt           string   //
	Confirm          string   //
	Message          string   // "..."
	Porcelain        bool     // true if `git status --porcelain` returns ""
}

// Init returns an initialized *Repo.
func Init(zw string, zu string, zr string, bp string, rn string) *Repo {

	r := new(Repo)

	r.BundlePath = tilde.Abs(bp) // "~/tmpgis" (bundle path)
	r.Workspace = zw             // "'main', 'go', 'bash' etc." (zone workspace)
	r.User = zu                  // "jychri" (zone user)
	r.Remote = zr                // "'github', 'gitlab'" (zone remote)
	r.Name = rn                  // "git-in-sync" (repo name)

	var b bytes.Buffer

	b.WriteString(r.BundlePath)
	if r.Workspace != "main" {
		b.WriteString("/")
		b.WriteString(r.Workspace)
	}
	r.WorkspacePath = b.String()
	// "/Users/jychri/tmpgis/go-lang/src/github.com/jychri"

	b.Reset()
	b.WriteString(r.WorkspacePath)
	b.WriteString("/")
	b.WriteString(r.Name)
	r.RepoPath = b.String()
	// "/Users/jychri/tmpgis/go-lang/src/github.com/jychri/git-in-sync"

	b.Reset()
	b.WriteString(r.RepoPath)
	b.WriteString("/.git")
	r.GitPath = b.String()
	// "/Users/jychri/tmpgis/go-lang/src/github.com/jychri/git-in-sync/.git"

	b.Reset()
	b.WriteString("--git-dir=")
	b.WriteString(r.GitPath)
	r.GitDir = b.String()
	// "--git-dir=/Users/jychri/tmpgis/go-lang/src/github.com/jychri/git-in-sync/.git"

	b.Reset()
	b.WriteString("--work-tree=")
	b.WriteString(r.RepoPath)
	r.WorkTree = b.String()
	// "--work-tree=/Users/jychri/tmpgis/go-lang/src/github.com/jychri/git-in-sync"

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
	// "https://github.com/jychri/git-in-sync"

	return r
}

// Error records the last error.
func (r *Repo) Error(dsc string, em string) {
	r.ErrorMessage = em
	r.ErrorName = dsc
	r.ErrorFirst = brf.First(em)

	if strings.Contains(r.ErrorFirst, "warning") {
		r.Verified = true
	}

	if strings.Contains(r.ErrorFirst, "fatal") {
		r.Verified = false
	}

	if strings.Contains(r.ErrorFirst, "Not a git repository") {
		r.Verified = false
	}

}

// VerifyWorkspace verifies that the r.WorkspacePath is present and accessible.
func (r *Repo) VerifyWorkspace(f flags.Flags, st *stat.Stat) {

	const dsc = "verify-workspace"

	if _, err := os.Stat(r.WorkspacePath); os.IsNotExist(err) {
		brf.Printv(f, "%v creating %v", emoji.Get("Folder"), r.WorkspacePath)
		os.MkdirAll(r.WorkspacePath, 0777)
		st.CreatedWorkspaces = append(st.CreatedWorkspaces, r.Workspace)
	}

	np := fchk.NoPermission(r.WorkspacePath)
	id := fchk.IsDirectory(r.WorkspacePath)

	switch {
	case id == true && np == false:
		r.Verified = true
		st.VerifiedWorkspaces = append(st.VerifiedWorkspaces, r.Workspace)
	case np == true:
		r.Error(dsc, "fatal: No permsission")
		st.InaccessibleWorkspaces = append(st.InaccessibleWorkspaces, r.Workspace)
	case id == false:
		r.Error(dsc, "fatal: No directory")
		st.InaccessibleWorkspaces = append(st.InaccessibleWorkspaces, r.Workspace)
	}
}

// GitSchedule ...
func (r *Repo) GitSchedule(f flags.Flags, st *stat.Stat) {

	const dsc = "GitSchedule"

	_, rerr := os.Stat(r.RepoPath)
	_, gerr := os.Stat(r.GitPath)

	switch {
	case fchk.IsFile(r.RepoPath):
		r.Error(dsc, "fatal: file occupying path")
	case fchk.IsDirectory(r.RepoPath) && fchk.NotEmpty(r.RepoPath) && os.IsNotExist(gerr):
		r.Error(dsc, "fatal: directory occupying path")
	case fchk.IsDirectory(r.RepoPath) && fchk.IsEmpty(r.RepoPath):
		r.PendingClone = true
		st.PendingClones = append(st.PendingClones, r.Name)
	case os.IsNotExist(rerr) && os.IsNotExist(gerr):
		r.PendingClone = true
		st.PendingClones = append(st.PendingClones, r.Name)
	case fchk.IsDirectory(r.RepoPath) && fchk.IsDirectory(r.GitPath):
		r.Verified = true
		r.PendingClone = false
	}
}

// GitClone clones a Git repository from r.URL if
// r.PendingClone is true.
func (r *Repo) GitClone(f flags.Flags) {
	const dsc = "GitClone"

	// return if !PendingClone
	if !r.PendingClone {
		return
	}

	// "cloning..."
	brf.Printv(f, "%v cloning %v {%v}", emoji.Get("Box"), r.Name, r.Workspace)

	args := []string{"clone", r.URL, r.RepoPath}
	if out, _ := r.git(args); out != "" {
		r.Error(dsc, out)
	} else {
		r.Cloned = true
	}
}

// GitConfigOriginURL gets the remote origin URL for a Repo.
func (r *Repo) GitConfigOriginURL() {
	const dsc = "GitConfigOriginURL"

	if !r.Verified {
		return
	}

	args := []string{r.GitDir, "config", "--get", "remote.origin.url"}
	out, _ := r.git(args)

	switch {
	case out == "":
		r.Error(dsc, "fatal: 'origin' does not appear to be a git repository")
	case out != r.URL:
		r.Error(dsc, "fatal: URL != OriginURL")
	default:
		r.OriginURL = out
	}
}

// GitRemoteUpdate ...
func (r *Repo) GitRemoteUpdate() {
	const dsc = "GitRemoteUpdate"

	if !r.Verified {
		return
	}

	args := []string{r.GitDir, r.WorkTree, "fetch", "origin"}
	_, em := r.git(args)

	// Warnings for redirects to "*./git" are ignored.
	wgit := strings.Join([]string{r.URL}, "/.git")

	switch {
	case strings.Contains(em, "warning: redirecting") && strings.Contains(em, wgit):
	case em != "":
		r.Error(dsc, em)
	}
}

// GitAbbrevRef ...
func (r *Repo) GitAbbrevRef() {
	const dsc = "GitAbbrevRef"

	if !r.Verified {
		return
	}

	args := []string{r.GitDir, r.WorkTree, "rev-parse", "--abbrev-ref", "HEAD"}
	if out, em := r.git(args); em != "" {
		r.Error(dsc, em)
	} else {
		r.LocalBranch = out
	}
}

// GitLocalSHA ...
func (r *Repo) GitLocalSHA() {
	const dsc = "GitLocalSHA"

	if !r.Verified {
		return
	}

	args := []string{r.GitDir, r.WorkTree, "rev-parse", "@"}
	if out, em := r.git(args); em != "" {
		r.Error(dsc, em)
	} else {
		r.LocalSHA = out
	}
}

// GitUpstreamBranch ...
func (r *Repo) GitUpstreamBranch() {
	const dsc = "GitUpstreamBranch"

	if !r.Verified {
		return
	}

	args := []string{r.GitDir, r.WorkTree, "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}"}
	if out, em := r.git(args); em != "" {
		r.Error(dsc, em)
	} else {
		r.UpstreamBranch = out
	}
}

// GitMergeBaseSHA ...
func (r *Repo) GitMergeBaseSHA() {
	const dsc = "GitMergeBaseSHA"

	if !r.Verified {
		return
	}

	args := []string{r.GitDir, r.WorkTree, "merge-base", "@", "@{u}"}
	if out, em := r.git(args); em != "" {
		r.Error(dsc, em)
	} else {
		r.MergeSHA = out
	}
}

// GitRevParseUpstream ...
func (r *Repo) GitRevParseUpstream() {
	const dsc = "GitRevParseUpstream"

	if !r.Verified {
		return
	}

	args := []string{r.GitDir, r.WorkTree, "rev-parse", "@{u}"}
	if out, em := r.git(args); em != "" {
		r.Error(dsc, em)
	} else {
		r.UpstreamSHA = out
	}
}

// GitDiffsNameOnly ...
func (r *Repo) GitDiffsNameOnly() {
	var out, em string
	const dsc = "GitDiffsNameOnly"

	if !r.Verified {
		return
	}

	args := []string{r.GitDir, r.WorkTree, "diff", "--name-only", "@{u}"}
	if out, em = r.git(args); em != "" {
		r.Error(dsc, em)
	}

	if out == "" {
		r.DiffsNameOnly = make([]string, 0)
		r.DiffsSummary = ""
	} else {
		r.DiffsNameOnly = strings.Fields(out)
		r.DiffsSummary = brf.Summary(r.DiffsNameOnly, 12)
	}
}

// GitShortstat ...
func (r *Repo) GitShortstat() {
	const dsc = "GitShortstat"

	if !r.Verified {
		return
	}

	// command
	args := []string{r.GitDir, r.WorkTree, "diff", "--shortstat"}
	if out, em := r.git(args); em != "" {
		r.Error(dsc, em)
	} else {
		r.ShortStat = out
	}

	// scrape with regular expressions
	// Set Insertions, Deletions
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

	// set Clean, ShortStatSummary
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

// GitUntracked ...
func (r *Repo) GitUntracked() {

	var out, em string
	const dsc = "GitUntracked"

	if !r.Verified {
		return
	}

	args := []string{r.GitDir, r.WorkTree, "ls-files", "--others", "--exclude-standard"}

	out, em = r.git(args)

	if em != "" {
		r.Error(dsc, em)
		return
	}

	if out == "" {
		r.UntrackedFiles = make([]string, 0)
		return
	}

	uf := strings.Fields(out)

	for _, f := range uf {
		f = path.Base(f)
		r.UntrackedFiles = append(r.UntrackedFiles, f)
	}

	r.UntrackedSummary = brf.Summary(r.UntrackedFiles, 12)
	r.Untracked = true
}

// SetStatus ...
func (r *Repo) SetStatus(f flags.Flags) {

	const dsc = "SetStatus"

	if !r.Verified {
		return
	}

	switch {
	case r.LocalSHA == "":
		r.Error(dsc, "fatal: r.LocalSHA = ''")
	case r.UpstreamSHA == "":
		r.Error(dsc, "fatal: r.UpstreamSHA = ''")
	case r.MergeSHA == "":
		r.Error(dsc, "fatal: r.MergeSHA = ''")
	case r.LocalSHA == r.UpstreamSHA:
		r.Status = "Up-To-Date"
	case r.LocalSHA == r.MergeSHA:
		r.Status = "Behind"
	case r.UpstreamSHA == r.MergeSHA:
		r.Status = "Ahead"
	}

	switch {
	case (r.Clean == true && r.Untracked == false && r.Status == "Ahead"):
		r.Category = "Pending"
		r.Status = "Ahead"
		r.Action = "Push"
	case (r.Clean == true && r.Untracked == false && r.Status == "Behind"):
		r.Category = "Pending"
		r.Status = "Behind"
		r.Action = "Pull"
	case (r.Clean == false && r.Untracked == false && r.Status == "Up-To-Date"):
		r.Category = "Pending"
		r.Status = "Dirty"
		r.Action = "Add-Commit-Push"
	case (r.Clean == false && r.Untracked == true && r.Status == "Up-To-Date"):
		r.Category = "Pending"
		r.Status = "DirtyUntracked"
		r.Action = "Add-Commit-Push"
	case (r.Clean == false && r.Untracked == false && r.Status == "Ahead"):
		r.Category = "Pending"
		r.Status = "DirtyAhead"
		r.Action = "Add-Commit-Push"
	case (r.Clean == false && r.Untracked == false && r.Status == "Behind"):
		r.Category = "Pending"
		r.Status = "DirtyBehind"
		r.Action = "Stash-Pull-Pop-Commit-Push"
	case (r.Clean == true && r.Untracked == true && r.Status == "Up-To-Date"):
		r.Category = "Pending"
		r.Status = "Untracked"
		r.Action = "Add-Commit-Push"
	case (r.Clean == true && r.Untracked == true && r.Status == "Ahead"):
		r.Category = "Pending"
		r.Status = "UntrackedAhead"
		r.Action = "Add-Commit-Push"
	case (r.Clean == false && r.Untracked == true && r.Status == "Behind"):
		r.Category = "Pending"
		r.Status = "UntrackedBehind"
	case (r.Clean == true && r.Untracked == false && r.Status == "Up-To-Date"):
		r.Category = "Complete"
		r.Status = "Up-To-Date"
	default:
		r.Category = "Skipped"
		r.Status = "Unknown"
		r.Error(dsc, "No matching conditions")
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
	switch {
	case f.Login() && r.Category == "Pending" && r.Status == "Behind":
		r.Category = "Scheduled"
	case f.Logout() && r.Category == "Pending" && r.Status == "Ahead":
		r.Category = "Scheduled"
	}

	var b bytes.Buffer
	var s string

	eb := emoji.Get("Bunny")
	rn := r.Name
	ub := r.UpstreamBranch
	et := emoji.Get("Turtle")
	ep := emoji.Get("Pig")
	dfs := len(r.DiffsNameOnly)
	ds := r.DiffsSummary
	sss := r.ShortStatSummary
	ufs := len(r.UntrackedFiles)
	us := r.UntrackedSummary

	switch r.Status {
	case "Ahead":
		s = fmt.Sprintf("%v %v is ahead of %v ", eb, rn, ub)
	case "Behind":
		s = fmt.Sprintf("%v %v is behind of %v ", et, rn, ub)
	case "Dirty", "DirtyUntracked", "DirtyAhead", "DirtyBehind":
		s = fmt.Sprintf("%v %v is dirty [%v]{%v}(%v)", ep, rn, dfs, ds, sss)
	case "Untracked", "UntrackedAhead", "UntrackedBehind":
		s = fmt.Sprintf("%v %v is untracked [%v]{%v}", ep, rn, ufs, us)
	}

	b.WriteString(s)
	s = ""

	switch r.Status {
	case "DirtyUntracked":
		s = fmt.Sprintf(" and untracked [%v]{%v}", ufs, us)
	case "DirtyAhead":
		s = fmt.Sprintf(" & ahead of %v", ub)
	case "DirtyBehind":
		s = fmt.Sprintf(" & behind  %v", ub)
	case "UntrackedAhead":
		s = fmt.Sprintf(" & is ahead of %v", ub)
	case "UntrackedBehind":
		s = fmt.Sprintf(" & is behind %v", ub)
	}

	if s != "" {
		b.WriteString(s)
	}

	r.Prompt = b.String()

	re := r.Remote
	er := emoji.Get("Rocket")
	eb = emoji.Get("Boat")
	ec := emoji.Get("Clipboard")

	switch r.Status {
	case "Ahead":
		s = fmt.Sprintf("%v push changes to %v? ", er, re)
	case "Behind":
		s = fmt.Sprintf("%v pull changes from %v? ", eb, re)
	case "Dirty":
		s = fmt.Sprintf("%v add all files, commit and push to %v? ", ec, re)
	case "DirtyUntracked":
		s = fmt.Sprintf("%v add all files, commit and push to %v? ", ec, re)
	case "DirtyAhead":
		s = fmt.Sprintf("%v add all files, commit and push to %v? ", ec, re)
	case "DirtyBehind":
		s = fmt.Sprintf("%v stash all files, pull changes, commit and push to %v? ", ec, re)
	case "Untracked":
		s = fmt.Sprintf("%v add all files, commit and push to %v? ", ec, re)
	case "UntrackedAhead":
		s = fmt.Sprintf("%v add all files, commit and push to %v? ", ec, re)
	case "UntrackedBehind":
		s = fmt.Sprintf("%v stash all files, pull changes, commit and push to %v? ", ec, re)
	}

	r.Confirm = s
}

// UserConfirm ...
func (r *Repo) UserConfirm(f flags.Flags) {

	if r.Category != "Pending" {
		return
	}

	if f.Mode == "oneline" {
		return
	}

	fmt.Println(r.Prompt)
	fmt.Printf(r.Confirm)

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

	// return if no commit message needed
	if ac := strings.Contains(r.Action, "Commit"); ac == false {
		return
	}

	em := emoji.Get("Memo")
	fmt.Printf("%v commit message: ", em)

	in, err = rdr.ReadString('\n')

	switch in {
	case "n", "no", "nah", "0", "stop", "skip", "abort", "halt", "quit", "exit", "":
		r.Category = "Skipped"
		r.Message = ""
	default:
		r.Category = "Scheduled"
		r.Message = in
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
	args := []string{"-C", r.RepoPath, "commit", "-m", r.Message}
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

	// return if !Verified
	if !r.Verified {
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
