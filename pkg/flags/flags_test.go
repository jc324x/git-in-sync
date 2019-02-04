package flags

import (
	"testing"

	"github.com/jychri/git-in-sync/pkg/atp"
)

func TestInit(t *testing.T) {
	f := Init()

	switch {
	case f.Config != "~/.gisrc.json":
		t.Errorf("Init: want: ~/.gisrc.json, got %v\n", f.Config)
	case f.Mode != "verify":
		t.Errorf("Init: want: verify, got %v\n", f.Mode)
	}
}

func TestTesting(t *testing.T) {
	for _, tr := range []struct {
		pkg, k string
	}{
		{"repos", "recipes"},
	} {

		p, cleanup := atp.Setup(tr.pkg, tr.k)

		defer cleanup()

		f := Testing(p)

		if p != f.Config {
			t.Errorf("Init: (want: %v, got: %v\n", p, f.Config)
		}
	}
}
