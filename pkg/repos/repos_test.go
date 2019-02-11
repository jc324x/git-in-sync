package repos

import (
	"os"
	"testing"

	"github.com/jychri/git-in-sync/pkg/atp"
	"github.com/jychri/git-in-sync/pkg/conf"
	"github.com/jychri/git-in-sync/pkg/flags"
	"github.com/jychri/git-in-sync/pkg/stat"
	"github.com/jychri/git-in-sync/pkg/timer"
)

func TestVerifyWorkspaces(t *testing.T) {

	for _, tr := range []struct {
		pkg, k string
	}{
		{"repos", "recipes"},
		{"repos", "google-apps-script"},
		{"repos", "tmp"},
	} {

		p, cleanup := atp.Setup(tr.pkg, tr.k)
		ti := timer.Init()
		f := flags.Testing(p)
		c := conf.Init(f)
		rs := Init(c, f, ti)
		st := stat.Init()

		defer cleanup()

		rs.VerifyWorkspaces(f, st, ti)

		for _, r := range rs {
			if _, err := os.Stat(r.WorkspacePath); os.IsNotExist(err) {
				t.Errorf("VerifyWorkspaces: %v does not exist", r.WorkspacePath)
			}
		}
	}
}

func TestVerifyRepos(t *testing.T) {

	for _, tr := range []struct {
		pkg, k string
	}{
		{"repos", "recipes"},
		{"repos", "google-apps-script"},
		{"repos", "tmp"},
	} {

		p, cleanup := atp.Setup(tr.pkg, tr.k)
		ti := timer.Init()
		f := flags.Testing(p)
		c := conf.Init(f)
		rs := Init(c, f, ti)
		st := stat.Init()

		defer cleanup()

		rs.VerifyWorkspaces(f, st, ti)
		rs.VerifyRepos(f, st, ti)

		// if _, err := os.Stat(tp); os.IsNotExist(err) {
		// 	t.Errorf("VerifyWorkspaces: %v does not exist", td)
		// }

	}
}
