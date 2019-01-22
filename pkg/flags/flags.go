package flags

import (
	"flag"
	"strings"
)

// Flags hold flag input
type Flags struct {
	Mode    string
	OneLine bool
	Count   int
	Summary string
}

func Init() (f Flags) {

	var ol bool  // one line
	var m string // mode
	var fc int   // flag count
	var s string // summary

	flag.StringVar(&m, "m", "verify", "mode")
	flag.BoolVar(&ol, "ol", false, "one-line")
	flag.Parse()

	// collect and join (e)nabled (f)lags
	var ef []string

	// mode
	if m != "" {
		fc += 1
	}

	// default value for mode is verify
	switch m {
	case "login", "logout", "verify":
	default:
		m = "verify"
	}
	ef = append(ef, m)

	// one-line
	if ol == true {
		fc += 1
		ef = append(ef, "one-line")
	}

	// summary
	s = strings.Join(ef, ", ")

	// set Flags
	f = Flags{m, ol, fc, s}

	return f
}

func (f Flags) Check(s string) bool {
	switch {
	case s == "one" && f.OneLine:
		return true
	case s == "logout" && f.Mode == "logout":
		return true
	case s == "login" && f.Mode == "login":
		return true
	default:
		return false
	}
}
