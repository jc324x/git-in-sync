package repo

import (
	"path"
	"strings"
	"testing"

	"github.com/jychri/git-in-sync/flags"
	"github.com/jychri/git-in-sync/tilde"
)

// private

func (r *Repo) eval(want string, dsc string, t *testing.T) {
	en := r.ErrorName
	em := r.ErrorMessage
	var m bool

	switch {
	case strings.Contains(em, "fatal:"):
		m = true
	case strings.Contains(em, "Not a git"):
		m = true
	case strings.Contains(em, "No matching"):
		m = true
	default:
		m = false
	}

	if !m {
		t.Errorf("Uncaught fatal error: %v %v", en, em)
	} else {
		r.Verified = true // true to keep testing
	}
}

func TestInit(t *testing.T) {
	zw := "main"
	zu := "jychri"
	zr := "github"
	bp := "~/tmpgis"
	rn := "git-in-sync"

	r := Init(zw, zu, zr, bp, rn)

	bp = tilde.Abs(bp)

	if r.BundlePath != bp {
		t.Errorf("Init: BundlePath got %v != want %v", r.BundlePath, bp)
	}

	if r.Workspace != zw {
		t.Errorf("Init: Workspace %v != %v", r.Workspace, zw)
	}

	if r.User != zu {
		t.Errorf("Init: User %v != %v", r.User, zu)
	}

	if r.Remote != zr {
		t.Errorf("Init: Remote %v != %v", r.Remote, zr)
	}

	if r.Name != rn {
		t.Errorf("Init: Name %v != %v", r.Name, rn)
	}

	twp := path.Join(bp)

	if r.WorkspacePath != twp {
		t.Errorf("Init: WorkspacePath %v != %v", r.WorkspacePath, twp)
	}

	trp := path.Join(bp, rn)

	if r.RepoPath != trp {
		t.Errorf("Init: RepoPath %v != '%v'", r.RepoPath, trp)
	}

	tgp := path.Join(bp, rn, ".git")

	if r.GitPath != tgp {
		t.Errorf("Init: GitPath %v != '%v'", r.GitPath, tgp)
	}

	tgd := path.Join(bp, rn, ".git")
	tgd = strings.Join([]string{"--git-dir=", tgd}, "")

	if r.GitDir != tgd {
		t.Errorf("Init: GitDir %v != '%v'", r.GitDir, tgd)
	}

	twt := path.Join(bp, rn)
	twt = strings.Join([]string{"--work-tree=", twt}, "")

	if r.WorkTree != twt {
		t.Errorf("Init: WorkTree %v != '%v'", r.WorkTree, twt)
	}

	zr = strings.Join([]string{zr, ".com"}, "")
	turl := strings.Join([]string{zr, zu, rn}, "/")
	turl = strings.Join([]string{"https://", turl}, "")

	if r.URL != turl {
		t.Errorf("Init: got %v != want %v", r.URL, turl)
	}
}

func TestErrors(t *testing.T) {

	zw := "fake"
	zu := "fake"
	zr := "github"
	bp := "~/fake"
	rn := "fake"

	r := Init(zw, zu, zr, bp, rn)
	r.Verified = true

	dsc := "GitConfigOriginURL"
	r.GitConfigOriginURL()
	want := "fatal: 'origin' does not appear to be a git repository"
	r.eval(want, dsc, t)

	dsc = "GitRemoteUpdate"
	r.GitRemoteUpdate()
	r.eval(want, dsc, t)

	dsc = "GitAbbrevRef"
	r.GitAbbrevRef()
	r.eval(want, dsc, t)

	dsc = "GitLocalSHA"
	r.GitLocalSHA()
	r.eval(want, dsc, t)

	dsc = "GitRevParseUpstream"
	r.GitRevParseUpstream()
	r.eval(want, dsc, t)

	dsc = "GitMergeBaseSHA"
	r.GitMergeBaseSHA()
	r.eval(want, dsc, t)

	dsc = "GitRevParseUpstream"
	r.GitMergeBaseSHA()
	r.eval(want, dsc, t)

	dsc = "GitDiffsNameOnly"
	r.GitDiffsNameOnly()
	r.eval(want, dsc, t)

	dsc = "GitShortstat"
	r.GitShortstat()
	r.eval(want, dsc, t)

	dsc = "GitUntracked"
	r.GitUntracked()
	r.eval(want, dsc, t)

	dsc = "SetStatus"
	f := flags.Testing("~/fakegisrc.json")
	r.SetStatus(f)
	r.eval(want, dsc, t)
}
