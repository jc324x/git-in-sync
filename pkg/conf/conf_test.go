package conf

import (
	"os"
	"reflect"
	"testing"

	"github.com/jychri/git-in-sync/pkg/atp"
	"github.com/jychri/git-in-sync/pkg/flags"
)

func TestInit(t *testing.T) {
	p := atp.Setup("conf", "recipes")

	f := flags.Flags{Mode: "verify", Config: p}

	c := Init(f)

	bs := c.Bundles[0]

	zs := bs.Zones

	// yeah, but with a function instead
	ts := atp.Tmap["recipes"]

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

	if err := os.Remove(p); err != nil {
		t.Errorf("Init (%v)", err.Error())
	}
}
