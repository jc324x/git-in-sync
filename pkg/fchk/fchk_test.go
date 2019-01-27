package fchk

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func GetTestDir() (os.FileInfo, string) {
	var abs string
	var err error
	var fi os.FileInfo

	if abs, err = filepath.Abs(""); err != nil {
		log.Fatalf("GetTestDir: Unable to access current directory")
	}

	s := path.Join(abs, "test_dir")

	if fi, err = os.Stat(s); err != nil {
		log.Fatalf("GetTestDir: Unable to get file info")
	}

	return fi, s
}

func GetTestFile() (os.FileInfo, string) {
	var abs string
	var err error
	var fi os.FileInfo

	if abs, err = filepath.Abs(""); err != nil {
		log.Fatalf("GetTestFile: Unable to access current directory")
	}

	s := path.Join(abs, "test_dir", "test_file")

	if fi, err = os.Stat(s); err != nil {
		log.Fatalf("GetTestDir: Unable to get file info")
	}

	return fi, s
}

func TestNoPermission(t *testing.T) {
	var want bool
	var fm os.FileMode

	tf, s := GetTestFile()

	want = true
	fm = 0444

	err := os.Chmod(s, fm)

	if err != nil {
		fmt.Println(err)
		t.Errorf("NoPermission: Chmod %v error (%v) ", fm, s)
	}

	got := NoPermission(tf)

	if got != want {
		t.Errorf("NoPermission: %v (%v != %v)", tf.Name(), got, want)
	}

	want = false
	fm = 0644

	err = os.Chmod(s, fm)

	if err != nil {
		fmt.Println(err)
		t.Errorf("NoPermission: Chmod %v error (%v) ", fm, s)
	}

	got = NoPermission(tf)

	if got != want {
		t.Errorf("NoPermission: %v (%v != %v)", tf.Name(), got, want)
	}
}
