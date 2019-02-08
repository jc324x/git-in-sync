package quit

import (
	"fmt"
	"log"
	"os"
)

// Active returns true if MODE != TESTING
func Active() bool {
	if ac := os.Getenv("MODE"); ac == "TESTING" {
		return false
	}
	return true
}

// Eval evaluates a condition. It MODE != TESTING, it returns a bool and a string
// values. If the condition is false and MODE != TESTING, the program will exit
// with a call to log.Fatalf(). A true condition with MODE != TESTING will return
// a bool and a string which will likely be ignored.
func Eval(b bool, fmats []string, v ...interface{}) (bool, string) {

	ac := Active()

	if len(fmats) != 2 {
		log.Fatalf("fmats != 2 : fmats[0] == false. fmats[1] == true")
	}

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
