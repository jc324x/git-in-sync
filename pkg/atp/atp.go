// Package atp (a test package) sets up a test environment for git-in-sync.
package atp

import (
	// "encoding/json"
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

// var tmap
var tmap = map[string][]struct {
	User, Remote, Workspace string
	Repos                   []string
}{}

// TSS is a test slice of structs.
var TSS = []struct {
	User, Remote, Workspace string
	Repos                   []string
}{
	{"hendricius", "github", "recipes", []string{"pizza-dough", "the-bread-code"}},
	{"cocktails-for-programmers", "github", "recipes", []string{"cocktails-for-programmers"}},
	{"rochacbruno", "github", "recipes", []string{"vegan_recipes"}},
	{"niw", "github", "recipes", []string{"ramen"}},
}

// Setup creates "~/tmpgis/%pkg/gisrc.json",
// writes the sample JSON to gisrc.json and
// returns the path to the file.
func Setup(pkg string, jkey string) string {
	if pkg == "" {
		log.Fatalf("pkg is empty")
	}

	if _, ok := jmap[jkey]; ok != true {
		log.Fatalf("%v not found in jmap", jkey)
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

	if err = ioutil.WriteFile(p, jmap[jkey], 0777); err != nil {
		log.Fatalf("Unable to write to %v (%v)", p, err.Error())
	}

	return p
}

// Clean removes "~/tmpgis/%pkg" and all child file/directories.
func Clean() {

}
