package fchk

import (
	"io"
	"os"
)

// NoPermission is a function that does a thing...
func NoPermission(info os.FileInfo) bool {

	if info == nil {
		return false
	}

	if len(info.Mode().String()) <= 4 {
		return true
	}

	if s := info.Mode().String()[1:4]; s != "rwx" {
		return true
	}

	return false
}

// IsDirectory is a function that does a thing...
func IsDirectory(info os.FileInfo) bool {
	if info == nil {
		return false
	}

	if info.IsDir() {
		return true
	}
	return false
}

// IsEmpty is a function that does a thing...
func IsEmpty(p string) bool {
	f, err := os.Open(p)

	if err != nil {
		return false
	}

	_, err = f.Readdir(1)

	if err == io.EOF {
		return true
	}

	return false
}

// NotEmpty is a function that does a thing...
func NotEmpty(p string) bool {
	f, err := os.Open(p)

	if err != nil {
		return false
	}

	_, err = f.Readdir(1)

	if err == io.EOF {
		return false
	}

	return true
}

// IsFile is a function that does a thing...
func IsFile(info os.FileInfo) bool {
	if info == nil {
		return false
	}

	if info.IsDir() {
		return false
	}
	return true
}
