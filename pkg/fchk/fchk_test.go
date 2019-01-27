package fchk

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func GetTestDir() (os.FileInfo, error) {
	var abs string
	var err error

	if abs, err = filepath.Abs(""); err != nil {
		log.Fatalf("GetTestDir: Unable to access current directory")
	}

	td := path.Join(abs, "test_dir")
	return os.Stat(td)
}

func TestNoPermission(t *testing.T) {
	want := false

	if td, err := GetTestDir(); err == nil {

		got := NoPermission(td)

		if got != want {
			t.Errorf("NoPermission: %v (%v != %v)", td.Name(), got, want)
		}
	}
}
