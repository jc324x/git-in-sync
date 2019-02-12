package repo

import (
	"path"
	"strings"
	"testing"

	"github.com/jychri/git-in-sync/pkg/tilde"
)

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
