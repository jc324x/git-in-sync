package repos

import (
// "os"
// "path"
// "testing"

// "github.com/jychri/git-in-sync/pkg/atp"
// "github.com/jychri/git-in-sync/pkg/conf"
// "github.com/jychri/git-in-sync/pkg/flags"
// "github.com/jychri/git-in-sync/pkg/timer"
)

// func TestVerifyWorkspaces(t *testing.T) {

// 	for _, tr := range []struct {
// 		pkg, k string
// 	}{
// 		{"repos", "recipes"},
// 	} {

// 		p, cleanup := atp.Setup(tr.pkg, tr.k)
// 		ti := timer.Init()
// 		f := flags.Testing(p)
// 		c := conf.Init(f)
// 		rs := Init(c, f, ti)

// 		defer cleanup()

// 		rs.VerifyWorkspaces(f, ti)

// 		td := atp.Dir(tr.pkg)
// 		tp := path.Join(td, tr.k)

// 		if _, err := os.Stat(tp); os.IsNotExist(err) {
// 			t.Errorf("VerifyWorkspaces: %v does not exist", td)
// 		}

// 	}
// }

// func TestVerifyCloned(t *testing.T) {

// 	for _, tr := range []struct {
// 		pkg, k string
// 	}{
// 		{"repos", "recipes"},
// 	} {

// 		p, cleanup := atp.Setup(tr.pkg, tr.k)
// 		ti := timer.Init()
// 		f := flags.Testing(p)
// 		c := conf.Init(f)
// 		rs := Init(c, f, ti)

// 		defer cleanup()

// 		rs.VerifyWorkspaces(f, ti)
// 		rs.VerifyCloned(f, ti)

// 		// td := atp.Dir(tr.pkg)
// 		// tp := path.Join(td, tr.k)

// 		// if _, err := os.Stat(tp); os.IsNotExist(err) {
// 		// 	t.Errorf("VerifyWorkspaces: %v does not exist", td)
// 		// }

// 	}
// }
