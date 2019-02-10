package main

import (
	"testing"
)

func TestInit(t *testing.T) {

	// do a thing here that creates a tmp ~/.gisrc.json or verifies it

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
