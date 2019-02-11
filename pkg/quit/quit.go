// Package quit wraps log.Fatalf() and adds additional options.
package quit

import (
	"fmt"
	"log"
	"os"
)

// private

func checkMode() bool {
	if env := os.Getenv("MODE"); env != "TESTING" {
		return true
	}
	return false
}

// Public

// Out ...
type Out struct {
	Status  bool
	Summary string
}

// Err evaluates an error, returning an Out struct
// if possible. If err != nil and mode != "TESTING",
// Err will call log.Fatalf.
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

// Bool evaluates an bool, returning an Out struct
// if possible. If bool == false and mode != "TESTING",
// Err will call log.Fatalf.
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
