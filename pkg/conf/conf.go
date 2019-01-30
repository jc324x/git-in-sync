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

// Init returns unmarshalled data from gisrc.json.
func Init(f flags.Flags) (c Config) {
	var u *user.User
	var err error
	var p string
	var bs []byte

	if u, err = user.Current(); err != nil {
		log.Fatalf("Can't id the current user (%v)", err.Error())
	}

	switch f.Config {
	case "~/.gisrc.json":
		p = fmt.Sprintf("%v/.gisrc.json", u.HomeDir)
	default:
		p = f.Config
	}

	if bs, err = ioutil.ReadFile(p); err != nil {
		log.Fatalf("Can't read file at (%v)\n", p)
	}

	if err = json.Unmarshal(bs, &c); err != nil {
		log.Fatalf("Can't unmarshal JSON from (%v)\n", p)
	}

	return c
}
