package conf

import (
	"reflect"
	"testing"

	"github.com/jychri/git-in-sync/pkg/atp"
	"github.com/jychri/git-in-sync/pkg/flags"
)

func TestInit(t *testing.T) {

	for _, tr := range []struct {
		pkg, recipe string
	}{
		{"conf", "recipes"},
		{"conf", "google-apps-script"},
		{"conf", "tmp"},
	} {

		p, clean := atp.Setup(tr.pkg, tr.recipe)

		defer clean()

		f := flags.Testing(p)

		c := Init(f)

		bs := c.Bundles[0]

		zs := bs.Zones

		rs := atp.Resulter(tr.recipe)

		for i := range rs {

			if rs[i].User != zs[i].User {
				t.Errorf("Init: (%v != %v)", rs[i].User, zs[i].User)
			}

			if rs[i].Remote != zs[i].Remote {
				t.Errorf("Init: (%v != %v)", rs[i].Remote, zs[i].Remote)
			}

			if rs[i].Workspace != zs[i].Workspace {
				t.Errorf("Init: (%v != %v)", rs[i].Workspace, zs[i].Workspace)
			}

			if !reflect.DeepEqual(rs[i].Repos, zs[i].Repos) {
				t.Errorf("Init: (%v != %v)", rs[i].Repos, zs[i].Repos)
			}

		}
	}
}
