// Package tilde supports "~/" as a shortcut for the current home directory.
package tilde

import (
	"log"
	"os/user"
	"path"
	"strings"
)

// AbsUser expands "~/" to "User/$user/" and returns a clean path.
// Given an absolute path, it returns a clean path.
func AbsUser(p string) string {

	var u *user.User

	u, err := user.Current()

	if err != nil {
		log.Fatalf("Unable to identify current user")
	}

	if !path.IsAbs(p) {
		return path.Join(u.HomeDir, strings.TrimPrefix(p, "~/"))
	}

	return path.Clean(p)
}
