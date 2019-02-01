package conf

import (
	"path"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/jychri/git-in-sync/pkg/flags"
)

func TestInit(t *testing.T) {
	var abs string
	var err error

	if abs, err = filepath.Abs(""); err != nil {
		t.Errorf("TestInit: err = %v\n", err.Error())
	}

	p := path.Join(abs, "ex_gisrc.json")

	f := flags.Flags{Mode: "verify", Config: p}

	c := Init(f)

	bs := c.Bundles[0]

	zs := bs.Zones

	ts := []struct {
		user, remote, workspace string
		repos                   []string
	}{
		{"hendricius", "github", "recipes", []string{"pizza-dough"}},
	}

	for i := 0; i < 1; i++ {

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
}
