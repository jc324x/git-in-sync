package fchk

import (
	"io"
	"os"
)

// NoPermission returns true target file is inaccessible.
func NoPermission(info os.FileInfo) bool {

	if info == nil {
		return false
	}

   if file, err := os.Open(path); err != nil {
           log.Print("NO => Cant open file for reading: ", path)
           return false


	if len(info.Mode().String()) <= 4 {
		return true
	}

	if s := info.Mode().String()[1:4]; s != "rwx" {
		return true
	}

	return false
}

// IsDirectory returns true if the target file is a directory.
func IsDirectory(info os.FileInfo) bool {
	if info == nil {
		return false
	}

	if info.IsDir() {
		return true
	}
	return false
}

// IsEmpty returns true if the target file is an empty directory.
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

// NotEmpty returns true if the target file is an non-empty directory.
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

// IsFile returns true if the target file is a file.
func IsFile(info os.FileInfo) bool {
	if info == nil {
		return false
	}

	if info.IsDir() {
		return false
	}
	return true
}
