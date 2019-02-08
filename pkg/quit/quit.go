package quit

import (
	"fmt"
	"log"
	"os"
)

// Quit ...
type Quit struct {
	B bool
	S string
}

// Active returns true if MODE != TESTING
func Active() bool {
	if ac := os.Getenv("MODE"); ac == "TESTING" {
		return false
	}
	return true
}

// Err ...
func Err(err error, fmats []string, v ...interface{}) Quit {
	var ac bool

	if env := os.Getenv("MODE"); env != "TESTING" {
		ac = true
	}

	switch {
	case ac && err == nil:
		return Quit{true, fmt.Sprintf(fmats[1], v...)}
	case ac && err != nil:
		log.Fatalf(fmats[0], v...)
	case !ac && err == nil:
		return Quit{true, fmt.Sprintf(fmats[1], v...)}
	case !ac && err != nil:
		return Quit{false, fmt.Sprintf(fmats[1], v...)}
	}

	return Quit{false, ""}
	// return false, ""
}

// Bool evaluates a boolean condition. It MODE != TESTING, it returns a bool and a string
// values. If the condition is false and MODE != TESTING, the program will exit
// with a call to log.Fatalf(). A true condition with MODE != TESTING will return
// a bool and a string which will likely be ignored.
func Bool(b bool, fmats []string, v ...interface{}) (bool, string) {
	var ac bool

	ac = Active()

	switch {
	case ac && b:
		return true, fmt.Sprintf(fmats[1], v...)
	case ac && !b:
		log.Fatalf(fmats[0], v...)
	case !ac && b:
		return true, fmt.Sprintf(fmats[1], v...)
	case !ac && !b:
		return false, fmt.Sprintf(fmats[0], v...)
	}

	return false, ""
}
