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
// and standard error messages as strings out and em.
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

// gitP runs a Git command and records the error
// values in r.ErrorName and r.ErrorMessage.
func (r *Repo) gitP(args []string, dsc string) {

	if r.Verified == false {
		return
	}

	var errb bytes.Buffer

	cmd := exec.Command("git", args...)
	cmd.Stderr = &errb
	cmd.Run()

	em := errb.String()
	em = strings.TrimSuffix(em, "\n")

	if em != "" {
		r.ErrorName = dsc
		r.ErrorMessage = em
	}
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
	Changed          int      // count of changed files (x)
	Insertions       int      // count of inserted files (y)
	Deletions        int      // count of deleted files (z)
	ShortStatSummary string   // "+y|-z" or "D" for Deleted if (x >= 1 && y == 0 && z == 0)
	Clean            bool     // true if Changed, Insertions and Deletions are all 0
	Untracked        bool     // true if if len(r.UntrackedFiles) >= 1
	UntrackedFiles   []string // `git ls-files --others --exclude-standard`, [a, b, c, d, e]
	UntrackedSummary string   // "a, b, c..."
	Category         string   // Complete, Pending, Skipped, Scheduled
	Status           string   // Complete is the last step
	Action           string   // Push, Pull, Add-Commit-Push etc.
	Prompt1          string   // First prompt message
	Prompt2          string   // Second prompt message
	Message          string   // Commit message
}

// Init returns an initialized *Repo.
func Init(workspace string, user string, remote string, bundle string, name string) *Repo {

	bundle = tilde.Abs(bundle) // set bundle to absolute path
	r := new(Repo)             // new Repo
	r.BundlePath = bundle      // ~/tmpgis etc.
	r.Workspace = workspace    // main, go, bash etc.
	r.User = user              // jychri
	r.Remote = remote          // github, gitlab etc.
	r.Name = name              // git-in-sync

	// /Users/jychri/tmpgis/golang or /Users/jychri/tmpgis (main)
	if workspace != "main" {
		r.WorkspacePath = path.Join(bundle, workspace)
	} else {
		r.WorkspacePath = bundle
	}

	// /Users/jychri/tmpgis/golang/src/github.com/jychri/git-in-sync
	r.RepoPath = path.Join(r.WorkspacePath, name)

	// /Users/jychri/tmpgis/go-lang/src/github.com/jychri/git-in-sync/.git
	r.GitPath = path.Join(r.RepoPath, ".git")

	// --git-dir=/Users/jychri/tmpgis/go-lang/src/github.com/jychri/git-in-sync/.git
	r.GitDir = strings.Join([]string{"--git-dir=", r.GitPath}, "")

	// --work-tree=/Users/jychri/tmpgis/go-lang/src/github.com/jychri/git-in-sync
	r.WorkTree = strings.Join([]string{"--work-tree=", r.RepoPath}, "")

	switch r.Remote {
	case "github":
		r.URL = strings.Join([]string{"https://github.com", r.User, r.Name}, "/")
	case "gitlab":
		r.URL = strings.Join([]string{"https://gitlab.com", r.User, r.Name}, "/")
	}

	return r
}

// Error records and evaluates last error.
func (r *Repo) Error(dsc string, em string) {
	r.ErrorMessage = em
	r.ErrorName = dsc

	switch {
	case strings.Contains(r.ErrorMessage, "warning"):
		r.Verified = true
	case strings.Contains(r.ErrorMessage, "fatal"):
		r.Verified = false
	case strings.Contains(r.ErrorMessage, "Not a git repository"):
		r.Verified = false
	}
}

// VerifyWorkspace verifies that the r.WorkspacePath is present and accessible.
func (r *Repo) VerifyWorkspace(f flags.Flags, st *stat.Stat) {

	const dsc = "VerifyWorkspace"

	if _, err := os.Stat(r.WorkspacePath); os.IsNotExist(err) {
		flags.Printv(f, "%v creating %v", emoji.Get("Folder"), r.WorkspacePath)
		// os.MkdirAll(r.WorkspacePath, 0444) // this breaks things, good for testing
		os.MkdirAll(r.WorkspacePath, 0766)
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

// GitClone clones a Git repository from r.URL.
func (r *Repo) GitClone(f flags.Flags) {
	const dsc = "GitClone"

	if !r.PendingClone {
		return
	}

	// "cloning..."
	flags.Printv(f, "%v cloning %v {%v}", emoji.Get("Box"), r.Name, r.Workspace)

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
	mod := strings.TrimSuffix(out, ".git")

	switch {
	case out == "":
		r.Error(dsc, "fatal: 'origin' does not appear to be a git repository")
	case mod == r.URL:
		r.OriginURL = out
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
	case strings.Contains(em, "fatal"):
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
		// log.Printf("%v: Shortstat ERR: %v | %v", r.Name, out, em)
	} else {
		r.ShortStat = out
		// log.Printf("%v: Shortstat OUT: %v", r.Name, out)
	}

	// scrape with regular expressions
	// Set Changed, Insertions, Deletions
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
		r.ShortStatSummary = "D"
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
		r.Status = "Complete"
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
	case (r.Clean == false && r.Untracked == false && r.Status == "Complete"):
		r.Category = "Pending"
		r.Status = "Dirty"
		r.Action = "Add-Commit-Push"
	case (r.Clean == false && r.Untracked == true && r.Status == "Complete"):
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
	case (r.Clean == true && r.Untracked == true && r.Status == "Complete"):
		r.Category = "Pending"
		r.Status = "Untracked"
		r.Action = "Add-Commit-Push"
	case (r.Clean == true && r.Untracked == true && r.Status == "Ahead"):
		r.Category = "Pending"
		r.Status = "UntrackedAhead"
		r.Action = "Add-Commit-Push"
	case (r.Clean == true && r.Untracked == true && r.Status == "Behind"):
		r.Category = "Pending"
		r.Status = "UntrackedBehind"
		r.Action = "Stash-Pull-Pop-Commit-Push"
	case (r.Clean == true && r.Untracked == false && r.Status == "Complete"):
		r.Category = "Complete"
		r.Status = "Complete"
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
		r.Action = "Pull"
	case f.Logout() && r.Category == "Pending" && r.Status == "Ahead":
		// flags.Printv(f, "%v", r.Name)
		r.Category = "Scheduled"
		r.Action = "Push"
	}

	var b bytes.Buffer
	var s string

	eb := emoji.Get("Bunny")
	rn := r.Name
	ub := r.UpstreamBranch
	et := emoji.Get("Turtle")
	ep := emoji.Get("Pig")
	dfc := len(r.DiffsNameOnly)
	ds := r.DiffsSummary
	sss := r.ShortStatSummary
	ufc := len(r.UntrackedFiles)
	us := r.UntrackedSummary

	switch r.Status {
	case "Ahead":
		s = fmt.Sprintf("%v %v is ahead of %v ", eb, rn, ub)
	case "Behind":
		s = fmt.Sprintf("%v %v is behind of %v ", et, rn, ub)
	case "Dirty", "DirtyUntracked", "DirtyAhead", "DirtyBehind":
		s = fmt.Sprintf("%v %v is dirty [%v]{%v}(%v)", ep, rn, dfc, ds, sss)
	case "Untracked", "UntrackedAhead", "UntrackedBehind":
		s = fmt.Sprintf("%v %v is untracked [%v]{%v}", ep, rn, ufc, us)
	}

	b.WriteString(s)
	s = ""

	switch r.Status {
	case "DirtyUntracked":
		s = fmt.Sprintf(" and untracked [%v]{%v}", ufc, us)
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

	r.Prompt1 = b.String()

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

	r.Prompt2 = s
}

// UserConfirm prompts the user with prompts Prompt1 and Prompt2
// and records the response.
func (r *Repo) UserConfirm(f flags.Flags) {

	if r.Category != "Pending" {
		return
	}

	if f.Mode == "oneline" {
		return
	}

	// prompt and read
	fmt.Println(r.Prompt1)
	fmt.Printf(r.Prompt2)
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

	// return if the user declined the commit
	if r.Category == "Skipped" {
		return
	}

	// prompt and read
	em := emoji.Get("Memo")               // Memo emoji
	fmt.Printf("%v commit message: ", em) // print
	in, err = rdr.ReadString('\n')        // read

	switch in {
	case "n", "no", "nah", "0", "stop", "skip", "abort", "halt", "quit", "exit", "":
		r.Category = "Skipped"
		r.Message = ""
	default:
		r.Category = "Scheduled"
		r.Message = in
	}
}

// GitAdd ...
func (r *Repo) GitAdd(f flags.Flags) {
	const dsc = "GitAdd"         // description
	eo := emoji.Get("Outbox")    // Outbox emoji
	rn := r.Name                 // repo name
	dfc := len(r.DiffsNameOnly)  // count: diff files
	ufc := len(r.UntrackedFiles) // count: untracked files
	ds := r.DiffsSummary         // summary: diffs
	us := r.UntrackedSummary     // summary: untracked
	sss := r.ShortStatSummary    // summary: (+/-)

	switch r.Status {
	case "Dirty", "DirtyUntracked", "DirtyAhead", "DirtyBehind":
		flags.Printv(f, "%v %v adding changes [%v]{%v}(%v)", eo, rn, dfc, ds, sss)
	case "Untracked", "UntrackedAhead", "UntrackedBehind":
		flags.Printv(f, "%v %v adding new files [%v]{%v}", eo, rn, ufc, us)
	}

	args := []string{"-C", r.RepoPath, "add", "-A"}
	r.gitP(args, dsc) // arguments and command
}

// GitCommit ...
func (r *Repo) GitCommit(f flags.Flags) {
	const dsc = "GitCommit"      // description
	ef := emoji.Get("Fire")      // Fire emoji
	rn := r.Name                 // repo name
	dfc := len(r.DiffsNameOnly)  // count: diff files
	ufc := len(r.UntrackedFiles) // count: untracked files
	ds := r.DiffsSummary         // summary: diffs
	us := r.UntrackedSummary     // summary: untracked
	sss := r.ShortStatSummary    // summary: (+/-)

	switch r.Status {
	case "Dirty", "DirtyUntracked", "DirtyAhead", "DirtyBehind":
		flags.Printv(f, "%v %v committing changes [%v]{%v}(%v)", ef, rn, dfc, ds, sss)
	case "Untracked", "UntrackedAhead", "UntrackedBehind":
		flags.Printv(f, "%v %v committing new files [%v]{%v}", ef, rn, ufc, us)
	}

	args := []string{"-C", r.RepoPath, "commit", "-m", r.Message}
	r.gitP(args, dsc) // arguments and command
}

// GitStash ...
func (r *Repo) GitStash(f flags.Flags) {
	const dsc = "GitStash"                             // description
	es := emoji.Get("Squirrel")                        // Squirrel emoji
	rn := r.Name                                       // repo name
	flags.Printv(f, "%v  %v stashing changes", es, rn) // print
	args := []string{"-C", r.RepoPath, "stash"}        // arguments
	r.gitP(args, dsc)                                  // command
}

// GitPop ...
func (r *Repo) GitPop(f flags.Flags) {
	const dsc = "GitPop"                               // description
	ep := emoji.Get("Popcorn")                         // Popcorn emoji
	rn := r.Name                                       // repo name
	flags.Printv(f, "%v  %v popping changes", ep, rn)  // print
	args := []string{"-C", r.RepoPath, "stash", "pop"} // arguments
	r.gitP(args, dsc)                                  // command
}

// GitPull ...
func (r *Repo) GitPull(f flags.Flags) {
	const dsc = "GitPull"                                                 // description
	es := emoji.Get("Ship")                                               // Ship emoji
	rn := r.Name                                                          // repo name
	ub := r.UpstreamBranch                                                // upstream branch
	rr := r.Remote                                                        // remote
	flags.Printv(f, "%v %v pulling changes from %v @ %v", es, rn, ub, rr) // print
	args := []string{"-C", r.RepoPath, "pull"}                            // arguments
	r.gitP(args, dsc)                                                     // command
}

// GitPush ...
func (r *Repo) GitPush(f flags.Flags) {
	const dsc = "GitPush"                                               // description
	er := emoji.Get("Rocket")                                           // Rocket emoji
	rn := r.Name                                                        // repo name
	ub := r.UpstreamBranch                                              // upstream branch
	rr := r.Remote                                                      // remote
	flags.Printv(f, "%v %v pushing changes to %v @ %v", er, rn, ub, rr) // print
	args := []string{"-C", r.RepoPath, "push"}                          // arguments
	r.gitP(args, dsc)                                                   // command
}

// GitClear ...
func (r *Repo) GitClear() {
	const dsc = "GitClean"
	r.OriginURL = ""
	r.LocalBranch = ""
	r.LocalSHA = ""
	r.UpstreamBranch = ""
	r.MergeSHA = ""
	r.UpstreamSHA = ""
	r.DiffsNameOnly = nil
	r.DiffsSummary = ""
	r.ShortStat = ""
	r.Changed = 0
	r.Insertions = 0
	r.Deletions = 0
	r.Clean = true
	r.ShortStatSummary = ""
	r.UntrackedFiles = nil
	r.UntrackedSummary = ""
	r.Untracked = false
	r.Status = ""
}
