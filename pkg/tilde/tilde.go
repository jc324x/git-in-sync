// Package tilde  "~/" as a shortcut for the current home directory.
package tilde

import (
	"log"
	"os/user"
	"path"
	"strings"
)

// Abs replaces "~/" with "/User/$user/" and returns a clean path.
func Abs(p string) string {

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
