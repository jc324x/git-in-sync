// Package conf implements access to gisrc.json files.
package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/user"

	"github.com/jychri/git-in-sync/pkg/flags"
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

// Path ...
func Path(f flags.Flags) string {
	var u *user.User
	var err error

	if u, err = user.Current(); err != nil {
		log.Fatalf("Can't id the current user (%v)", err.Error())
	}

	switch f.Config {
	case "~/.gisrc.json":
		return fmt.Sprintf("%v/.gisrc.json", u.HomeDir)
	default:
		return f.Config
	}
}

// Init returns unmarshalled data from gisrc.json.
func Init(f flags.Flags) (c Config) {
	var bs []byte
	var err error
	var p string

	p = Path(f)

	if bs, err = ioutil.ReadFile(p); err != nil {
		log.Fatalf("Can't read file at (%v)\n", p)
	}

	if err = json.Unmarshal(bs, &c); err != nil {
		log.Fatalf("Can't unmarshal JSON from (%v)\n", p)
	}

	return c
}
