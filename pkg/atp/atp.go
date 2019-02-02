// Package atp (a test package) sets up a test environment for git-in-sync.
package atp

import (
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
)

var jmap = map[string][]byte{
	"recipes": []byte(`
		{
			"bundles": [{
				"path": "~/tmpgis",
				"zones": [{
						"user": "hendricius",
						"remote": "github",
						"workspace": "recipes",
						"repositories": [
							"pizza-dough",
							"the-bread-code"
						]
					},
					{
						"user": "cocktails-for-programmers",
						"remote": "github",
						"workspace": "recipes",
						"repositories": [
							"cocktails-for-programmers"
						]
					},
					{
						"user": "rochacbruno",
						"remote": "github",
						"workspace": "recipes",
						"repositories": [
							"vegan_recipes"
						]
					},
					{
						"user": "niw",
						"remote": "github",
						"workspace": "recipes",
						"repositories": [
							"ramen"
						]
					}
				]
			}]
		}
`),
}

// Setup creates "~/tmpgis/%pkg/gisrc.json",
// writes the sample JSON to gisrc.json and
// returns the path to the file.
func Setup(pkg string, k string) string {
	if pkg == "" {
		log.Fatalf("pkg is empty")
	}

	if _, ok := jmap[k]; ok != true {
		log.Fatalf("%v not found in jmap", k)
	}

	var u *user.User

	u, err := user.Current()

	if err != nil {
		log.Fatalf("Unable to identify current user (%v)", err.Error())
	}

	p := path.Join(u.HomeDir, "tmpgis", pkg)

	if err = os.MkdirAll(p, 0777); err != nil {
		log.Fatalf("Unable to create %v", p)
	}

	p = path.Join(p, "gisrc.json")

	if err = ioutil.WriteFile(p, jmap[k], 0777); err != nil {
		log.Fatalf("Unable to write to %v (%v)", p, err.Error())
	}

	return p
}

// Result ...
type Result struct {
	User, Remote, Workspace string
	Repos                   []string
}

// Results ...
type Results []Result

var rmap = map[string]Results{
	"recipes": {
		{"hendricius", "github", "recipes", []string{"pizza-dough", "the-bread-code"}},
		{"cocktails-for-programmers", "github", "recipes", []string{"cocktails-for-programmers"}},
		{"rochacbruno", "github", "recipes", []string{"vegan_recipes"}},
		{"niw", "github", "recipes", []string{"ramen"}},
	},
}

// Resulter returns expected results for testing.

// Wanter returns expected results for testing.
func Wanter(k string) Results {
	if _, ok := rmap[k]; ok != true {
		log.Fatalf("%v not found in rmap", k)
	}

	return rmap[k]
}

// Clean removes "~/tmpgis/%pkg" and all child file/directories.
func Clean() {

}
