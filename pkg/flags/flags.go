// Package flags adds Mode and Config flags to git-in-sync.
package flags

import (
	"flag"
	"log"
	"os/user"
	"path"
	"strings"
)

// Flags records values for Mode and Config.
type Flags struct {
	Mode   string
	Config string
}

// Init returns validated user input as Flags.
func Init() (f Flags) {

	var c, m string

	flag.StringVar(&m, "m", "verify", "mode")
	flag.StringVar(&c, "c", "~/.gisrc.json", "configuration")
	flag.Parse()

	switch m {
	case "login", "logout", "verify", "oneline":
	default:
		m = "verify"
	}

	var u *user.User

	u, err := user.Current()

	if err != nil {
		log.Fatalf("Unable to identify current user")
	}

	if !path.IsAbs(c) {
		c = path.Join(u.HomeDir, strings.TrimPrefix(c, "~/"))
	}

	c = path.Clean(c)

	return Flags{Mode: m, Config: c}
}

// Testing returns a Flags instance:
// Mode: testing
// Config: c
func Testing(c string) Flags {
	return Flags{Mode: "testing", Config: c}
}
