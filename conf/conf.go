// Package conf implements access to gisrc.json config files.
package conf

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/jychri/git-in-sync/flags"
)

// read returns the contents of the file at path f.Config as []byte.
// If something goes wrong, stop the run.
func read(f flags.Flags) (bs []byte) {
	var err error

	if bs, err = ioutil.ReadFile(f.Config); err != nil {
		log.Fatalf("Unable to read file at %v \n (%v)", f.Config, err)
	}

	return bs
}

// unmarsmall unmarshalls []byte into Config
// If something goes wrong, stop the run.
func unmarshal(bs []byte) (c Config) {
	if err := json.Unmarshal(bs, &c); err != nil {
		log.Fatalf("Unable to unmarshal data\n")
	}

	return c
}

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
func Init(f flags.Flags) (c Config) {
	bs := read(f)     // read the file at path f.Config
	c = unmarshal(bs) // unmarshal the data from file at path f.Config
	return c
}
