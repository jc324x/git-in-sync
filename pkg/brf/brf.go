// Package brf implements functions with brief output.
package brf

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os/user"
	"path"
	"strings"

	"github.com/jychri/git-in-sync/pkg/flags"
)

// Printv calls prints to standard output if not running in 'oneline' or 'testing' mode.
func Printv(f flags.Flags, s string, z ...interface{}) (err error) {

	if f.Mode != "oneline" && f.Mode != "testing" {
		fmt.Println(fmt.Sprintf(s, z...))
		return
	}

	return errors.New("N/A")
}

// Single returns a string slice with no duplicate entries.
func Single(ssl []string) (sl []string) {

	smap := make(map[string]bool)

	for i := range ssl {
		if smap[ssl[i]] == true {
		} else {
			smap[ssl[i]] = true
			sl = append(sl, ssl[i])
		}
	}

	return sl
}

// Summary returns a set length string summarizing the contents of a string slice.
func Summary(sl []string, l int) string {
	if len(sl) == 0 {
		return ""
	}

	var csl []string // check slice
	var b bytes.Buffer

	for _, s := range sl {
		lc := len(strings.Join(csl, ", ")) // (l)ength(c)heck
		switch {
		case lc <= l-10 && len(s) <= 20: //
			csl = append(csl, s)
		case lc <= l && len(s) <= 12:
			csl = append(csl, s)
		}
	}

	b.WriteString(strings.Join(csl, ", "))

	if len(sl) != len(csl) {
		b.WriteString("...")
	}

	return b.String()
}

// First returns the first line from a multi line string.
func First(s string) string {
	lines := strings.Split(strings.TrimSuffix(s, "\n"), "\n")

	if len(lines) >= 1 {
		return lines[0]
	}

	return ""
}

// AbsUser expands "~/" to "User/$user/" and returns a clean path.
// Given an absolute path, it returns a clean path.
func AbsUser(s string) string {

	var u *user.User

	u, err := user.Current()

	if err != nil {
		log.Fatalf("Unable to identify current user")
	}

	if !path.IsAbs(s) {
		return path.Join(u.HomeDir, strings.TrimPrefix(s, "~/"))
	}

	return path.Clean(s)
}
