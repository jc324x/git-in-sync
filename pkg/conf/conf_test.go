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

	ts := []struct {
		user, remote, workspace string
		repos                   []string
	}{
		{"hendricius", "github", "recipes", []string{"pizza-dough", "the-bread-code"}},
		{"cocktails-for-programmers", "github", "recipes", []string{"cocktails-for-programmers"}},
		{"rochacbruno", "github", "recipes", []string{"vegan_recipes"}},
		{"niw", "github", "recipes", []string{"ramen"}},
	}

	for i := range ts {

		if ts[i].user != zs[i].User {
			t.Errorf("Init: (%v != %v)", ts[i].user, zs[i].User)
		}

		if ts[i].remote != zs[i].Remote {
			t.Errorf("Init: (%v != %v)", ts[i].remote, zs[i].Remote)
		}

		if ts[i].workspace != zs[i].Workspace {
			t.Errorf("Init: (%v != %v)", ts[i].workspace, zs[i].Workspace)
		}

		if !reflect.DeepEqual(ts[i].repos, zs[i].Repos) {
			t.Errorf("Init: (%v != %v)", ts[i].repos, zs[i].Repos)
		}

	}

	if err = os.Remove(p); err != nil {
		t.Errorf("Init (%v)", err.Error())
	}
}
