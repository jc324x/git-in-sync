package flags

import (
	"testing"
)

func TestInit(t *testing.T) {
	f := Init()

	switch {
	case f.Config != "~/.gisrc.json":
		t.Errorf("TestInit: want: ~/.gisrc.json, got %v\n", f.Config)
	case f.Mode != "verify":
		t.Errorf("TestInit: want: verify, got %v\n", f.Mode)
	}
}
