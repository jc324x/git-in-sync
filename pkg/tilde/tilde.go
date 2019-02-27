// Package tilde  "~/" as a shortcut for the current home directory.
package tilde

import (
	"log"
	"os/user"
	p "path"
	"strings"
)

// Abs replaces "~/" with "/User/$user/" and returns a clean path.
func Abs(path string) string {

	var u *user.User

	u, err := user.Current()

	if err != nil {
		log.Fatalf("Unable to identify current user")
	}

	if !p.IsAbs(path) {
		return p.Join(u.HomeDir, strings.TrimPrefix(path, "~/"))
	}

	return p.Clean(path)
}
