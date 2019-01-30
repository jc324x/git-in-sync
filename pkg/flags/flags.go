package flags

import (
	"flag"
	"strings"
)

// Flags hold flag input
type Flags struct {
	Mode    string
	Config  string
	Count   int
	Summary string
}

// Init ...
func Init() (f Flags) {

	var c, m, s string // config, mode, summary
	var fc int         // flag count

	flag.StringVar(&m, "m", "verify", "mode")
	flag.StringVar(&c, "c", "~/.gisrc.json", "configuration")
	flag.Parse()

	// collect and join (e)nabled (f)lags
	var ef []string

	// mode
	if m != "" {
		fc++
	}

	ef = append(ef, c)
	fc++

	// set unsupported modes to "verify"
	switch m {
	case "login", "logout", "verify", "oneline":
	default:
		m = "verify"
	}
	ef = append(ef, m)

	// summary
	s = strings.Join(ef, ", ")

	// set Flags
	f = Flags{m, c, fc, s}

	return f
}
