package fchk

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"testing"
)

const (
	itmp string = "tmp"
	itde string = "empty"
	itda string = "active"
	itf  string = "file"
)

func makeTmp() string {
	var abs, p string
	var err error

	if p, err = filepath.Abs(""); err != nil {
		log.Fatalf("makeTmp (%v)", err.Error())
	}

	p = path.Join(abs, itmp)

	if err = os.Mkdir(p, 0755); err != nil {
		log.Fatalf("makeTmp (%v)", err.Error())
	}

	return p
}

func makeTmpDE() string {
	var abs, p string
	var err error

	if p, err = filepath.Abs(""); err != nil {
		log.Fatalf("makeTmpDE (%v)", err.Error())
	}

	p = path.Join(abs, itmp, itde)

	if err = os.Mkdir(p, 0755); err != nil {
		log.Fatalf("makeTmpDE (%v)", err.Error())
	}

	return p
}

func makeTmpDA() string {
	var abs, p string
	var err error

	if p, err = filepath.Abs(""); err != nil {
		log.Fatalf("makeTmpDA (%v)", err.Error())
	}

	p = path.Join(abs, itmp, itda)

	if err = os.Mkdir(p, 0755); err != nil {
		log.Fatalf("makeTmpDA (%v)", err.Error())
	}

	return p
}

func makeTmpF() string {
	var p string
	var err error

	if p, err = filepath.Abs(""); err != nil {
		log.Fatalf("makeTmpF (%v)", err.Error())
	}

	p = path.Join(p, itmp, itda, itf)

	if _, err = os.OpenFile(p, os.O_RDONLY|os.O_CREATE, 0755); err != nil {
		log.Fatalf("makeTmpF (%v)", err.Error())
	}

	return p
}

var tmp = makeTmp()   // .../fchk/tmp (dir)
var tda = makeTmpDA() // .../fchk/tmp/empty (dir)
var tde = makeTmpDE() // .../fchk/tmp/active (dir)
var tf = makeTmpF()   // .../fchk/tmp/active/file (file)

func TestNoPermission(t *testing.T) {

	for _, c := range []struct {
		in   string
		want bool
	}{
		{tmp, false},
		{tda, false},
		{tde, false},
		{tf, false},
	} {

		got := NoPermission(c.in)

		if got != c.want {
			t.Errorf("NoPermission: (got: %v,  want: %v)", got, c.want)
		}
	}
}

func TestIsDir(t *testing.T) {

	for _, c := range []struct {
		in   string
		want bool
	}{
		{tmp, true},
		{tda, true},
		{tde, true},
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
	for _, c := range []struct {
		in   string
		want bool
	}{
		{tmp, false},
		{tda, false},
		{tde, true},
		{tf, false},
	} {

		got, err := IsEmpty(c.in)

		if err != nil {
			t.Errorf("IsEmpty: err = %v\n", err.Error())
		}

		if got != c.want {
			t.Errorf("IsEmpty: (got: %v, want: %v) {%v}\n", got, c.want, c.in)
		}
	}
}

func TestClean(t *testing.T) {
	var p string
	var err error

	if p, err = filepath.Abs(""); err != nil {
		t.Errorf("Clean: (%v)", err.Error())
	}

	p = path.Join(p, tmp)

	if err = os.RemoveAll(p); err != nil {
		t.Errorf("Clean: (%v)", err.Error())
	}
}
