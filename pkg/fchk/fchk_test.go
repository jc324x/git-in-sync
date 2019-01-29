package fchk

import (
	// "io/ioutil"
	"log"
	"os"
	"path"
	// "path/filepath"
	"testing"
)

func makeTmpD() (string, string) {
	tmp := os.TempDir()
	ed := path.Join(tmp, "fchk_test", "empty_dir")

	if err := os.MkdirAll(ed, 0666); err != nil {
		log.Fatalf("makeTmpD: Error creating '.../fchk_test/empty_dir' (%v)", err.Error())
	}

	ad := path.Join(tmp, "fchk_test", "active_dir")

	if err := os.MkdirAll(ad, 0666); err != nil {
		log.Fatalf("makeTmpD: Error creatin  '.../fchk_test/active_dir' (%v)", err.Error())
	}

	return ed, ad
}

func makeTmpF(ad string) string {
	p := path.Join(ad, "test_file")

	if _, err := os.Create(p); err != nil {
		log.Fatalf("makeTmpF: Error creating '.../fchk_test/active_dir/test_file'")
	}

	return p
}

var ad, ed = makeTmpD()
var f = makeTmpF(ad)

func TestNoPermissionP(t *testing.T) {

	for _, c := range []struct {
		want bool
		fm   os.FileMode
	}{
		{false, 0666},
		{true, 0444},
	} {

		err := os.Chmod(f, c.fm)

		if err != nil {
			t.Errorf("NoPermission: Chmod %v error (%v) ", c.fm, f)
		}

		got := NoPermission(f)

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

func TestNotEmpty(t *testing.T) {
	t.Log(os.TempDir())
}

// func Testy(t *testing.T) {

// }
