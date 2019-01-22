package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/user"
)

// Config holds the data from ~/.gisrc.json after Unmasrhalling.
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

// initConfig returns data from ~/.gisrc.json as a Config struct.
func Init() (c Config) {

	// get the current user, otherwise fatal
	u, err := user.Current()

	if err != nil {
		log.Fatal(err)
	}

	// expand "~/" to "/Users/user"
	g := fmt.Sprintf("%v/.gisrc.json", u.HomeDir)

	// print
	// targetPrintln(f, "%v reading %v", e.Glasses, g)

	// read file
	r, err := ioutil.ReadFile(g)

	if err != nil {
		log.Fatalf("No file found at %v\n", g)
	}

	// unmarshall json
	err = json.Unmarshal(r, &c)

	if err != nil {
		log.Fatalf("Can't unmarshal JSON from %v\n", g)
	}

	return c
}
