// Package quit wraps log.Fatalf() and adds additional options.
package quit

import (
	"fmt"
	"log"
	"os"
)

// private

// checkMode returns true if the environment variable
// "MODE" == "TESTING".
func checkMode() bool {
	if env := os.Getenv("MODE"); env != "TESTING" {
		return true
	}
	return false
}

// Public

// Out holds Status and Summary values.
type Out struct {
	Status  bool
	Summary string
}

// Err evaluates an error, and returns Out,
// if environment variable "MODE" == "TESTING".
// If mode != "TESTING" && err != nil, execution
// will stop with log.Fatalf().
func Err(err error, fm []string, v ...interface{}) Out {

	var m = checkMode()

	switch {
	case m && err != nil:
		log.Fatalf(fm[0], v...)
	case !m && err != nil:
		return Out{false, fmt.Sprintf(fm[0], v...)}
	}

	return Out{true, fmt.Sprintf(fm[1], v...)}
}

// Bool evaluates a bool, and returns Out
// if environment variable "MODE" == "TESTING".
// If mode != "TESTING" && b != true, execution
// will stop with log.Fatalf().
func Bool(b bool, fm []string, v ...interface{}) Out {

	var m = checkMode()

	switch {
	case m && !b:
		log.Fatalf(fm[0], v...)
	case !m && !b:
		return Out{false, fmt.Sprintf(fm[0], v...)}
	}

	return Out{true, fmt.Sprintf(fm[1], v...)}
}
