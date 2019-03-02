package conf

import (
	"reflect"
	"testing"

	"github.com/jychri/git-in-sync/atp"
	"github.com/jychri/git-in-sync/flags"
)

func TestInit(t *testing.T) {

	for _, tr := range []struct {
		pkg, recipe string
	}{
		{"conf", "recipes"},
		{"conf", "tmpgis"},
	} {
		p, cleanup := atp.Setup(tr.pkg, tr.recipe)
		f := flags.Testing(p)
		c := Init(f)
		bs := c.Bundles[0]
		zs := bs.Zones
		rs := atp.Resulter(tr.recipe)

		defer cleanup()

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
