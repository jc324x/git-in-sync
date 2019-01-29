package fchk

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"testing"
)

const (
	tmpd string = "tmp_dir"
	tmpf string = "tmp_file"
)

func makeTmpD() string {
	var p string
	var err error

	if p, err = filepath.Abs(""); err != nil {
		log.Fatalf("makeTmpD (%v)", err.Error())
	}

	p = path.Join(p, tmpd)

	if err = os.Mkdir(p, 0777); err != nil {
		log.Fatalf("makeTmpD (%v)", err.Error())
	}

	return p
}

func makeTmpF() string {
	var p string
	var err error

	if p, err = filepath.Abs(""); err != nil {
		log.Fatalf("makeTmpF (%v)", err.Error())
	}

	p = path.Join(p, tmpd, tmpf)

	if _, err = os.OpenFile(p, os.O_RDONLY|os.O_CREATE, 0777); err != nil {
		log.Fatalf("makeTmpF (%v)", err.Error())
	}

	return p
}

var td = makeTmpD()
var tf = makeTmpF()

func TestNoPermissionP(t *testing.T) {

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

// func TestIsDir(t *testing.T) {
// 	// td := GetTestDir()
// 	// tf := GetTestFile()

// 	for _, c := range []struct {
// 		in   string
// 		want bool
// 	}{
// 		{td, true},
// 		{tf, false},
// 	} {

// 		got, err := IsDirectory(c.in)

// 		if err != nil {
// 			t.Errorf("IsDirectory: err = %v\n", err.Error())
// 		}

// 		if got != c.want {
// 			t.Errorf("IsDirectory: (got: %v, want: %v) {%v}\n", got, c.want, c.in)
// 		}
// 	}
// }

// func TestIsEmpty(t *testing.T) {
// 	td := GetTestDir()

// 	if IsEmpty(td) == true {

// 	}

// }

func TestClean(t *testing.T) {
	var p string
	var err error

	if p, err = filepath.Abs(""); err != nil {
		t.Errorf("Clean: (%v)", err.Error())
	}

	p = path.Join(p, tmpd)

	if err = os.RemoveAll(p); err != nil {
		t.Errorf("Clean: (%v)", err.Error())
	}
}
