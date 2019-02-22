// Package flags add Mode and Config flags to git-in-sync.
package flags

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

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

	flag.StringVar(&m, "mode", "verify", "mode")
	flag.StringVar(&c, "config", "~/.gisrc.json", "configuration")
	flag.Parse()

	switch m {
	case "login", "logout", "verify", "oneline", "testing":
	default:
		m = "verify"
	}

	if env := os.Getenv("MODE"); env == "TESTING" {
		m = "testing"
	}

	if env := os.Getenv("GISRC"); env != "" {
		c = env
	}

	c = tilde.Abs(c)

	return Flags{Mode: m, Config: c}
}

// Testing returns a Flags instance with Mode == "testing".
func Testing(c string) Flags {
	return Flags{Mode: "testing", Config: c}
}

// ClearScreen clears the screen.
func (f Flags) ClearScreen() {
	switch f.Mode {
	case "oneline":
	case "testing":
	default:
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

// Printv calls prints to standard output if not running in 'oneline' or 'testing' mode.
func Printv(f Flags, s string, z ...interface{}) string {
	out := fmt.Sprintf(s, z...)

	switch f.Mode {
	case "oneline":
	case "testing":
	default:
		fmt.Println(out)
	}
	return out
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
