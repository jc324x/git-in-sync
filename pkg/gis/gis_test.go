package main

import (
	"os"
	"testing"

	"github.com/jychri/git-in-sync/pkg/atp"
)

func TestInit(t *testing.T) {

	os.Setenv("MODE", "TESTING")
	_, cleanup := atp.Direct("gis", "recipes")

	defer cleanup()

	f, rs, _, ti := Init()

	if f.Config == "" {
		t.Errorf("Init: %v = ''", f.Config)
	}

	if f.Mode == "" {
		t.Errorf("Init: %v = ''", f.Mode)
	}

	if len(rs) <= 0 {
		t.Errorf("No repos %v", len(rs))
	}

	if _, err := ti.Get("Start"); err != nil {
		t.Errorf("No start moment in %+v", ti)
	}
}
