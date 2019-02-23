// Package fchk tests files and directories.
package fchk

import (
	"io"
	"io/ioutil"
	"os"
	"path"
)

// NoPermission returns true if the target can't be read,
// can't be written to or doesn't exist.
func NoPermission(p string) bool {

	var f *os.File
	var err error

	defer f.Close()

	if f, err = os.Open(p); err != nil {
		return true
	}

	if _, err := os.Stat(p); err != nil {
		return true
	}

	if f, err = os.OpenFile(p, os.O_WRONLY, 0766); err != nil {
		if os.IsPermission(err) {
			defer f.Close()
			return true
		}
	}

	if f, err = os.OpenFile(p, os.O_RDONLY, 0766); err != nil {
		if os.IsPermission(err) {
			return true
		}
	}

	filename := path.Join(p, "tmp")

	defer os.Remove(filename)

	if err = ioutil.WriteFile(filename, []byte{}, 0777); err != nil {
		return true
	}

	return false
}

// IsDirectory returns true if the target is a directory.
func IsDirectory(p string) bool {

	var fi os.FileInfo
	var err error

	if fi, err = os.Stat(p); err != nil {
		return false
	}

	if fi.IsDir() {
		return true
	}

	return false
}

// IsEmpty returns true if the target is an empty directory.
func IsEmpty(p string) bool {

	var f *os.File
	var err error

	if c := IsDirectory(p); c != true || err != nil {
		return false
	}

	f, err = os.Open(p)
	defer f.Close()

	if err != nil {
		return false
	}

	_, err = f.Readdir(1)

	if err == io.EOF {
		return true
	}

	return false
}

// NotEmpty returns true if the target is an non-empty directory.
func NotEmpty(p string) bool {

	var f *os.File
	var err error

	if f, err = os.Open(p); err != nil {
		return false
	}

	if _, err = f.Readdir(1); err == io.EOF {
		return false
	}

	return true
}

// IsFile returns true if the target is a file.
func IsFile(p string) bool {

	var fi os.FileInfo
	var err error

	if fi, err = os.Stat(p); err != nil {
		return false
	}

	if fi.IsDir() {
		return false
	}

	return true
}
