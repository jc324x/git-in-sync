package fchk

import (
	"fmt"
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

func TestSet(t *testing.T) {
	tf := GetTestFile()
	var m os.FileMode

	m = 0666
	os.Chmod(tf, m)
	fmt.Printf("0666 -> %v\n", NoPermissionP(tf))

	m = 0444
	os.Chmod(tf, m)
	fmt.Printf("0444 -> %v\n", NoPermissionP(tf))
}

func TestNoPermissionP(t *testing.T) {
	tf := GetTestFile()

	for i, c := range []struct {
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

		got := NoPermissionP(tf)

		if got != c.want {
			t.Errorf("NoPermission: (%v != %v) #%v", got, c.want, i)
		}
	}
}
