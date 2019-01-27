package fchk

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func GetTestDir() string {
	var abs string
	var err error

	if abs, err = filepath.Abs(""); err != nil {
		log.Fatalf("GetTestDir: Unable to access current directory")
	}

	return path.Join(abs, "test_dir")
}

func GetTestFile() string {
	var abs string
	var err error

	if abs, err = filepath.Abs(""); err != nil {
		log.Fatalf("GetTestFile: Unable to access current directory")
	}

	return path.Join(abs, "test_dir", "test_file")
}

func TestNoPermissionP(t *testing.T) {
	tf := GetTestFile()

	for _, c := range []struct {
		want bool
		fm   os.FileMode
	}{
		{false, 0666},
		{true, 0444},
	} {

		err := os.Chmod(tf, c.fm)

		if err != nil {
			t.Errorf("NoPermission: Chmod %v error (%v) ", c.fm, tf)
		}

		got := NoPermission(tf)

		if got != c.want {
			t.Errorf("NoPermission: (got: %v,  want: %v)", got, c.want)
		}
	}
}

func TestIsDir(t *testing.T) {
	td := GetTestDir()
	tf := GetTestFile()

	for _, c := range []struct {
		in   string
		want bool
	}{
		{td, true},
		{tf, false},
	} {

		got, err := IsDirectory(c.in)

		if err != nil {
			t.Errorf("IsDirectory: err = %v\n", err.Error())
		}

		if got != c.want {
			t.Errorf("IsDirectory: (got: %v, want: %v) {%v}\n", got, c.want, c.in)
		}
	}
}

func TestIsEmpty(t *testing.T) {
	td := GetTestDir()

	if IsEmpty(td) == true {

	}

}

func TestNotEmpty(t *testing.T) {
	// td := GetTestDir()

}
