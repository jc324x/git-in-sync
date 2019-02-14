package repo

import (
	"path"
	"strings"
	"testing"

	"github.com/jychri/git-in-sync/pkg/flags"
	"github.com/jychri/git-in-sync/pkg/tilde"
)

// private

// erm returns Repo values ErrorName and
// ErrorMessage as strings. It's used by
// TestErrors to streamline the test pattern.
func (r *Repo) erm() (string, string) {
	return r.ErrorName, r.ErrorMessage
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
		t.Errorf("Init: BundlePath %v != %v", r.BundlePath, bp)
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
		t.Errorf("Init: %v != %v", r.URL, turl)
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

	var en, em string

	dsc := "GitConfigOriginURL"
	r.GitConfigOriginURL()
	want := "fatal: 'origin' does not appear to be a git repository"

	if en, em = r.erm(); en != dsc || em != want {
		t.Errorf("%v %v != %v", dsc, r.ErrorMessage, want)
	}

	if r.Verified == true {
		t.Errorf("verified != false (%v) ", dsc)
	} else {
		r.Verified = true
	}

	dsc = "GitRemoteUpdate"
	r.GitRemoteUpdate()

	if en, em = r.erm(); en != dsc && em == "" {
		t.Errorf("%v em = '%v'", dsc, r.ErrorMessage)
	}

	if r.Verified == true {
		t.Errorf("verified != false (%v) ", dsc)
	} else {
		r.Verified = true
	}

	dsc = "GitAbbrevRef"
	r.GitAbbrevRef()

	if en, em = r.erm(); en != dsc && em == "" {
		t.Errorf("%v em = '%v'", dsc, r.ErrorMessage)
	}

	if r.Verified == true {
		t.Errorf("verified != false (%v) ", dsc)
	} else {
		r.Verified = true
	}

	dsc = "GitLocalSHA"
	r.GitLocalSHA()

	if en, em = r.erm(); en != dsc && em == "" {
		t.Errorf("%v em = '%v'", dsc, r.ErrorMessage)
	}

	if r.Verified == true {
		t.Errorf("verified != false (%v) ", dsc)
	} else {
		r.Verified = true
	}

	dsc = "GitRevParseUpstream"
	r.GitRevParseUpstream()

	if en, em = r.erm(); en != dsc && em == "" {
		t.Errorf("%v em = '%v'", dsc, r.ErrorMessage)
	}

	if r.Verified == true {
		t.Errorf("verified != false (%v) ", dsc)
	} else {
		r.Verified = true
	}

	dsc = "GitMergeBaseSHA"
	r.GitMergeBaseSHA()

	if en, em = r.erm(); en != dsc && em == "" {
		t.Errorf("%v em = '%v'", dsc, r.ErrorMessage)
	}

	if r.Verified == true {
		t.Errorf("verified != false (%v) ", dsc)
	} else {
		r.Verified = true
	}

	dsc = "GitRevParseUpstream"
	r.GitMergeBaseSHA()

	if en, em = r.erm(); en != dsc && em == "" {
		t.Errorf("%v em = '%v'", dsc, r.ErrorMessage)
	}

	if r.Verified == true {
		t.Errorf("verified != false (%v) ", dsc)
	} else {
		r.Verified = true
	}

	dsc = "GitDiffsNameOnly"
	r.GitDiffsNameOnly()

	if en, em = r.erm(); en != dsc && em == "" {
		t.Errorf("%v em = '%v'", dsc, r.ErrorMessage)
	}

	if r.Verified == true {
		t.Errorf("verified != false (%v) ", dsc)
	} else {
		r.Verified = true
	}

	dsc = "GitShortstat"
	r.GitShortstat()

	if en, em = r.erm(); en != dsc && em == "" {
		t.Errorf("%v em = '%v'", dsc, r.ErrorMessage)
	}

	if r.Verified == true {
		t.Errorf("verified != false (%v) ", dsc)
	} else {
		r.Verified = true
	}

	dsc = "GitUntracked"
	r.GitUntracked()

	if en, em = r.erm(); en != dsc && em == "" {
		t.Errorf("%v em = '%v'", dsc, r.ErrorMessage)
	}

	if r.Verified == true {
		t.Errorf("verified != false (%v) ", dsc)
	} else {
		r.Verified = true
	}

	dsc = "SetStatus"
	f := flags.Testing("~/fakegisrc.json")
	r.SetStatus(f)

	if en, em = r.erm(); en != dsc && em == "" {
		t.Errorf("%v em = '%v'", dsc, r.ErrorMessage)
	}

	if r.Verified == true {
		t.Errorf("verified != false (%v) ", dsc)
	} else {
		r.Verified = true
	}

}
