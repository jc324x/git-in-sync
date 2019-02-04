// Package fchk tests files and directories.
package fchk

import (
	"io"
	"os"
)

// NoPermission returns true if the target can't be read,
// can't be written to or doesn't exist.
func NoPermission(p string) (bool, error) {

	var f *os.File
	var err error

	if f, err = os.Open(p); err != nil {
		return true, err
	}

	defer f.Close()

	if _, err = os.Stat(p); err != nil {
		return true, err
	}

	if f, err = os.OpenFile(p, os.O_WRONLY, 0666); err != nil {
		if os.IsPermission(err) {
			return true, err
		}
	}

	defer f.Close()

	if f, err = os.OpenFile(p, os.O_RDONLY, 0666); err != nil {
		if os.IsPermission(err) {
			return true, err
		}
	}

	defer f.Close()

	return false, nil
}

// IsDirectory returns true if the target is a directory.
func IsDirectory(p string) (bool, error) {

	var fi os.FileInfo
	var err error

	if fi, err = os.Stat(p); err != nil {
		return false, err
	}

	if fi.IsDir() {
		return true, nil
	}

	return false, nil
}

// IsEmpty returns true if the target is an empty directory.
func IsEmpty(p string) (bool, error) {

	var f *os.File
	var err error

	if c, err := IsDirectory(p); c != true || err != nil {
		return false, err
	}

	f, err = os.Open(p)
	defer f.Close()

	if err != nil {
		return false, err
	}

	_, err = f.Readdir(1)

	if err == io.EOF {
		return true, nil
	}

	return false, nil
}

// NotEmpty returns true if the target is an non-empty directory.
func NotEmpty(p string) (bool, error) {

	var f *os.File
	var err error

	if f, err = os.Open(p); err != nil {
		return false, err
	}

	if _, err = f.Readdir(1); err == io.EOF {
		return false, err
	}

	return true, nil
}

// IsFile returns true if the target is a file.
func IsFile(p string) (bool, error) {

	var fi os.FileInfo
	var err error

	if fi, err = os.Stat(p); err != nil {
		return false, err
	}

	if fi.IsDir() {
		return false, nil
	}

	return true, nil
}
