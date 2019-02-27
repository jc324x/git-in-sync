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

// make ./tmp directory
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

// make ./tmp/empty
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

// make ./tmp/active
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

// make ./tmp/active/file
func makeTmpF() string {
	var p string
	var err error

	if p, err = filepath.Abs(""); err != nil {
		log.Fatalf("makeTmpF (%v)", err.Error())
	}

	p = path.Join(p, itmp, itda, itf)

	if _, err = os.OpenFile(p, os.O_RDONLY|os.O_CREATE, 0766); err != nil {
		log.Fatalf("makeTmpF (%v)", err.Error())
	}

	return p
}

var tmp = makeTmp()   // .../fchk/tmp (dir)
var tda = makeTmpDA() // .../fchk/tmp/empty (dir)
var tde = makeTmpDE() // .../fchk/tmp/active (dir)
var tf = makeTmpF()   // .../fchk/tmp/active/file (file)

func TestNoPermission(t *testing.T) {

	for _, tr := range []struct {
		in   string
		want bool
	}{
		{tmp, false},
		{tda, false},
		{tde, false},
		{tf, false},
	} {

		got := NoPermission(tr.in)

		if got != tr.want {
			t.Errorf("NoPermission (%v) got: %v !=  want: %v", tr.in, got, tr.want)
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

		got := IsDirectory(c.in)

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

		got := IsEmpty(c.in)

		if got != c.want {
			t.Errorf("IsEmpty: (got: %v, want: %v) {%v}\n", got, c.want, c.in)
		}
	}
}

func TestIsFile(t *testing.T) {
	for _, c := range []struct {
		in   string
		want bool
	}{
		{tmp, false},
		{tda, false},
		{tde, false},
		{tf, true},
	} {

		got := IsFile(c.in)

		if got != c.want {
			t.Errorf("IsFile: (got: %v, want: %v) {%v}\n", got, c.want, c.in)
		}
	}
}

// Removes the temporary directory and files created for these tests.
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
