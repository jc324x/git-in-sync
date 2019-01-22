package flags

import (
	"flag"
	"strings"
)

// Flags hold flag input
type Flags struct {
	Mode    string
	Count   int
	Summary string
}

func Init() (f Flags) {

	var c, m, s string // config, mode, summary
	var fc int         // flag count

	flag.StringVar(&c, "c", "~/.gisrc.json", "configuration")
	flag.StringVar(&m, "m", "verify", "mode")
	flag.Parse()

	// collect and join (e)nabled (f)lags
	var ef []string

	// mode
	if m != "" {
		fc += 1
	}

	ef = append(ef, c)
	fc += 1

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
	f = Flags{m, fc, s}

	return f
}

func (f Flags) Check(s string) bool {
	switch f.Mode {
	case "oneLine":
		return true
	// case s == "logout" && f.Mode == "logout":
	// 	return true
	// case s == "login" && f.Mode == "login":
	// 	return true
	default:
		return false
	}
}
