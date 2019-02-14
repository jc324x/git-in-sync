// Package flags add Mode and Config flags to git-in-sync.
package flags

import (
	"flag"

	"github.com/jychri/git-in-sync/pkg/tilde"
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
	case "login", "logout", "verify", "oneline", "testing":
	default:
		m = "verify"
	}

	c = tilde.Abs(c)

	return Flags{Mode: m, Config: c}
}

// Testing returns a Flags instance with Mode == "testing".
func Testing(c string) Flags {
	return Flags{Mode: "testing", Config: c}
}

// Login returns true if f.Mode == "login".
func (f Flags) Login() bool {
	if f.Mode == "login" {
		return true
	}
	return false
}

// Logout returns true if f.Mode == "logout".
func (f Flags) Logout() bool {
	if f.Mode == "logout" {
		return true
	}
	return false

}
