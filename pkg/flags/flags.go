package flags

import (
	"flag"
)

// Flags ...
type Flags struct {
	Mode   string
	Config string
}

// Init ...
func Init() (f Flags) {

	var c, m string

	flag.StringVar(&m, "m", "verify", "mode")
	flag.StringVar(&c, "c", "~/.gisrc.json", "configuration")
	flag.Parse()

	// set unsupported modes to "verify"
	switch m {
	case "login", "logout", "verify", "oneline":
	default:
		m = "verify"
	}

	// set Flags
	f = Flags{Mode: m, Config: c}

	return f
}
