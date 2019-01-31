package conf

import (
	"fmt"
	"path"
	"path/filepath"
	"testing"

	"github.com/jychri/git-in-sync/pkg/flags"
)

func TestInit(t *testing.T) {
	var abs, p string
	var err error
	var c Config

	if abs, err = filepath.Abs(""); err != nil {
		t.Errorf("TestInit: err = %v\n", err.Error())
	}

	p = path.Join(abs, "ex_gisrc.json")

	var f = flags.Flags{Mode: "verify", Config: p}

	if c = Init(f); len(c.Bundles) != 1 {
		t.Errorf("TestInit: err = bundle mismatch")
	}

	b := c.Bundles[0]

	for _, z := range b.Zones {
		fmt.Println(z.Remote)
	}

	// for i := range []struct {
	// 	user, remote, workspace string
	// 	repos                   []string
	// }{
	// 	{"hendricius", "github", "recipes", []string{"pizza-dough"}},
	// } {
	// 	fmt.Println(i)
	// 	// fmt.Printf("%v", zs[i])

	// 	// if x.user != zs[i].User {
	// 	// 	t.Errorf("shit")
	// 	// }

	// 	// if x.remote != zs[i].Remote {
	// 	// 	t.Errorf("shit")
	// 	// }

	// 	// if x.user != zs[i].User {
	// 	// 	t.Errorf("shit")
	// 	// }

	// }
}
