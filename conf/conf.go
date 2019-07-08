// Package conf implements access to gisrc.json files.
package conf

import (
	"encoding/json"
	"io/ioutil"

	"github.com/jychri/git-in-sync/flags"
)

// private

// read reads the file at f.Config. read returns
// a byte slice and quit.Out. In a normal run,
// quit.Err will call log.Fatalf() if err != nil.
// In a test run, quit.Err will record and return
// the failure as quit.Out.
func read(f flags.Flags) (bs []byte) {
	var err error

	if bs, err = ioutil.ReadFile(f.Config); err != nil {
		// handle error
	}

	return bs
}

// unmarsmall unmarhalls the contents of a gisrc.json file,
// read to a byte slice by read, and returns a Config and
// quit.Out. In a normal run quit.Err will call log.Fatalf()
// if err != nil. In a test run, quit.Err will record and
// return the failure as quit.Out.
func unmarshal(bs []byte, f flags.Flags) (c Config) {
	if err := json.Unmarshal(bs, &c); err != nil {
		// handle error
	}

	return c
}

// Public

// Config holds unmrashalled JSON from a gisrc.json file.
type Config struct {
	Bundles []struct {
		Path  string `json:"path"`
		Zones []struct {
			User      string   `json:"user"`
			Remote    string   `json:"remote"`
			Workspace string   `json:"workspace"`
			Repos     []string `json:"repositories"`
		} `json:"zones"`
	} `json:"bundles"`
}

// Init returns unmarshalled data from gisrc.json.
// The Flags' Mode and Config values are validated
// prior to their use here.
func Init(f flags.Flags) (c Config) {
	bs := read(f)        // read the file at path f.Config
	c = unmarshal(bs, f) // unmarshal the data from file at path f.Config
	return c
}
