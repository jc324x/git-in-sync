package repos

import (
	"os"
	"strings"
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
		{"repos-workspaces", "recipes"},
	} {
		// p, cleanup := atp.Setup(tr.pkg, tr.k)
		p, _ := atp.Setup(tr.pkg, tr.k)
		ti := timer.Init()
		f := flags.Testing(p)
		c := conf.Init(f)
		st := stat.Init()
		rs := Init(c, f, st, ti)

		// defer cleanup()

		rs.VerifyWorkspaces(f, st, ti)

		for _, r := range rs {
			if _, err := os.Stat(r.WorkspacePath); os.IsNotExist(err) {
				t.Errorf("VerifyWorkspaces: %v is missing", r.WorkspacePath)
			}
		}
	}
}

func TestVerifyRepos(t *testing.T) {

	for _, tr := range []struct {
		pkg, k string
	}{
		{"repos-repos", "recipes"},
	} {
		// p, cleanup := atp.Setup(tr.pkg, tr.k)
		p, _ := atp.Setup(tr.pkg, tr.k)
		ti := timer.Init()
		f := flags.Testing(p)
		c := conf.Init(f)
		st := stat.Init()
		rs := Init(c, f, st, ti)

		// defer cleanup()

		rs.VerifyWorkspaces(f, st, ti)
		rs.VerifyRepos(f, st, ti)

		for _, r := range rs {
			if _, err := os.Stat(r.GitPath); os.IsNotExist(err) {
				t.Errorf("#VerifyRepos %v is missing", r.GitPath)
			}

			if r.ErrorName != "" || r.ErrorMessage != "" {
				t.Errorf("#VerifyRepos (%v -> %v) @%v", r.ErrorName, r.ErrorMessage, r.Name)
			}
		}
	}
}

func TestVerifyChanges(t *testing.T) {

	for _, tr := range []struct {
		scope, k string
	}{
		{"repos-changes", "tmpgis"},
	} {
		// p, cleanup := atp.Hub(tr.pkg, tr.k)
		p, _ := atp.Hub(tr.scope, tr.k)
		ti := timer.Init()
		f := flags.Testing(p)
		c := conf.Init(f)
		st := stat.Init()
		rs := Init(c, f, st, ti)

		// defer cleanup()

		rs.VerifyWorkspaces(f, st, ti)
		rs.VerifyRepos(f, st, ti)

		// setup
		for _, r := range rs {

			if trim := strings.TrimPrefix(r.Name, "gis-"); trim != r.Status {
				t.Errorf("VerifyChanges: %v mismatch: %v != %v", r.Name, trim, r.Status)
			}

			r.Category = "Scheduled"
			r.Message = "'TESTVERIFY' commit"
		}

		rs.changesAsync(f, st, ti)
		rs.infoAsync(f, ti)

		for _, r := range rs {
			if r.Status != "Complete" {
				t.Errorf("VerifyChanges: %v %v != Complete", r.Name, r.Status)
			}
		}
	}
}
