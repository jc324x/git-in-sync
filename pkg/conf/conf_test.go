package conf

import (
	"reflect"
	"testing"

	"github.com/jychri/git-in-sync/pkg/atp"
	"github.com/jychri/git-in-sync/pkg/flags"
)

func TestInit(t *testing.T) {
	p, clean := atp.Setup("conf", "recipes", "gisrc")

	defer clean()

	f := flags.Flags{Mode: "verify", Config: p}

	c := Init(f)

	bs := c.Bundles[0]

	zs := bs.Zones

	ts := atp.Want("recipes")

	for i := range ts {

		if ts[i].User != zs[i].User {
			t.Errorf("Init: (%v != %v)", ts[i].User, zs[i].User)
		}

		if ts[i].Remote != zs[i].Remote {
			t.Errorf("Init: (%v != %v)", ts[i].Remote, zs[i].Remote)
		}

		if ts[i].Workspace != zs[i].Workspace {
			t.Errorf("Init: (%v != %v)", ts[i].Workspace, zs[i].Workspace)
		}

		if !reflect.DeepEqual(ts[i].Repos, zs[i].Repos) {
			t.Errorf("Init: (%v != %v)", ts[i].Repos, zs[i].Repos)
		}

	}
}
