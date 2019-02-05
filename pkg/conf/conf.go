// Package conf implements access to gisrc.json files.
package conf

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/jychri/git-in-sync/pkg/flags"
	// "github.com/jychri/git-in-sync/pkg/tilde"
)

// Config holds unmrashalled data from gisrc.json.
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
	var bs []byte
	var err error

	// if f.Config == "~/.gisrc.json" {
	// 	f.Config = brf.AbsUser("~/.gisrc.json")
	// }

	if bs, err = ioutil.ReadFile(f.Config); err != nil {
		log.Fatalf("Can't read file at (%v)\n", f.Config)
	}

	if err = json.Unmarshal(bs, &c); err != nil {
		log.Fatalf("Can't unmarshal JSON from (%v)\n", f.Config)
	}

	return c
}
