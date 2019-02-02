package flags

import (
	"flag"
)

// Flags holds the user settings for the current run.
type Flags struct {
	Mode   string
	Config string
}

// Init returns Flags with a validated mode and a default or set configuration.
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

	return Flags{Mode: m, Config: c}
}

// Testing is used to set values directly.
func Testing(c string) Flags {
	return Flags{Mode: "testing", Config: c}
}
