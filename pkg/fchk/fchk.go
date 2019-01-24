package fchk

import (
	"io"
	"log"
	"os"
	"os/user"
	// "path" -> Clean
	"strings"
)

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

func IsDirectory(info os.FileInfo) bool {
	if info == nil {
		return false
	}

	if info.IsDir() {
		return true
	} else {
		return false
	}
}

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

func IsFile(info os.FileInfo) bool {
	if info == nil {
		return false
	}

	if info.IsDir() {
		return false
	} else {
		return true
	}
}

func validatePath(p string) string {
	if t := strings.TrimPrefix(p, "~/"); t != p {
		u, err := user.Current()

		if err != nil {
			log.Fatalf("Unable to identify the current user")
		}

		t := strings.Join([]string{u.HomeDir, "/", t}, "")
		return strings.TrimSuffix(t, "/")
	}
	return strings.TrimSuffix(p, "/")
}

func lastPathSelection(p string) string {
	if strings.Contains(p, "/") == true {
		sp := strings.SplitAfter(p, "/") // split path
		lp := sp[len(sp)-1]              // last path
		return lp
	} else {
		return p
	}
}
