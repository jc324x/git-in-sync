// Package conf implements access to gisrc.json files.
package conf

import (
	"encoding/json"
	"io/ioutil"

	"github.com/jychri/goku/quitter"

	"github.com/jychri/git-in-sync/flags"
)

// private

// read reads the file at f.Config. read returns
// a byte slice and quit.Out. In a normal run,
// quit.Err will call log.Fatalf() if err != nil.
// In a test run, quit.Err will record and return
// the failure as quit.Out.
func read(f flags.Flags) ([]byte, quit.Out) {
	var bs []byte
	var err error

	bs, err = ioutil.ReadFile(f.Config)

	fm := []string{"Can't read file at (%v)\n", "Read file at (%v)\n"}
	qo := quit.Err(err, fm, f.Config)

	return bs, qo
}

// unmarsmall unmarhalls the contents of a gisrc.json file,
// read to a byte slice by read, and returns a Config and
// quit.Out. In a normal run quit.Err will call log.Fatalf()
// if err != nil. In a test run, quit.Err will record and
// return the failure as quit.Out.
func unmarshal(bs []byte, f flags.Flags) (Config, quit.Out) {
	var c Config
	var err error

	err = json.Unmarshal(bs, &c)

	fm := []string{"Can't unmarshal JSON from (%v)\n", "Unmarshalled (%v)\n"}
	qo := quit.Err(err, fm, f.Config)

	return c, qo
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
	bs, _ := read(f)        // read the file at path f.Config
	c, _ = unmarshal(bs, f) // unmarshal the data from file at path f.Config
	return c
}
