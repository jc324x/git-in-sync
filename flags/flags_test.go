package flags

import (
	"os"
	"testing"

	"github.com/jychri/git-in-sync/tilde"
)

func TestAll(t *testing.T) {

	os.Setenv("MODE", "TESTING")
	os.Setenv("GISRC", "~/.gisrc.json")

	f := Init()
	gj := tilde.Abs("~/.gisrc.json")

	switch {
	case f.Config != gj:
		t.Errorf("Flags: want: ~/.gisrc.json, got %v\n", f.Config)
	case f.Mode != "testing":
		t.Errorf("Flags: want: testing, got %v\n", f.Mode)
	}

	f.Mode = "login"

	if b := f.Login(); b != true {
		t.Errorf("Flags: want: true, got %v\n", b)
	}

	f.Mode = "logout"

	if b := f.Logout(); b != true {
		t.Errorf("Flags: want: true, got %v\n", b)
	}

	f.Mode = "testing"

	want := "This is a test."
	got := Printv(f, "This is a %v.", "test")

	if got != want {
		t.Errorf("Flags: want: %v, got %v\n", got, want)
	}
}
