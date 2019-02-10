// Package atp manages test environments for git-in-sync packages.
package atp

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
)

// private

var jmap = map[string][]byte{
	"recipes": []byte(`
		{
			"bundles": [{
				"path": "SETPATH",
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
`), "google-apps-script": []byte(`
		{
			"bundles": [{
				"path": "SETPATH",
				"zones": [{
						"user": "jychri",
						"remote": "github",
						"workspace": "google-apps-script",
						"repositories": [
							"crunchy-calendar",
							"daily-sign-up",
							"data-flipper",
							"easy-csv",
							"frequency-responder",
							"google-apps-script-cheat-sheet",
							"mega-merge",
							"missing-homework"
						]
					}
				]
			}]
		}`),
	"tmp": []byte(`
	{
		"bundles": [{
			"path": "SETPATH",
			"zones": [{
					"user": "jychri",
					"remote": "github",
					"workspace": "tmp",
					"repositories": [
						"tmp1",
						"tmp2",
						"tmp3",
						"tmp4",
						"tmp5"
					]
				}
			]
		}]
	}
	`),
}

var rmap = map[string]Results{
	"recipes": {
		{"hendricius", "github", "recipes", []string{"pizza-dough", "the-bread-code"}},
		{"cocktails-for-programmers", "github", "recipes", []string{"cocktails-for-programmers"}},
		{"rochacbruno", "github", "recipes", []string{"vegan_recipes"}},
		{"niw", "github", "recipes", []string{"ramen"}},
	},
	"google-apps-script": {
		{"jychri", "github", "google-apps-script", []string{
			"crunchy-calendar",
			"daily-sign-up",
			"data-flipper",
			"easy-csv",
			"frequency-responder",
			"google-apps-script-cheat-sheet",
			"mega-merge",
			"missing-homework",
		}}},
	"tmp": {
		{"jychri", "github", "google-apps-script", []string{
			"tmp1",
			"tmp2",
			"tmp3",
			"tmp4",
			"tmp5",
		}}},
}

// Public

// Setup creates a test environment at ~/tmpgis/$pkg/.
// ~/tmpgis/$pkg/ and ~/tmpgis/$pkg/gisrc.json are created,
// key k is matched to Jmap, returning j ([]byte) if valid,
// $td replaces 'SETPATH' in j, which is written to gisrc.json.
// Setup returns the absolute path of ~/tmpgis/$pkg/gisrc.json
// and a cleanup function that removes ~/tmpgis/$pkg/.
func Setup(pkg string, k string) (string, func()) {

	var j []byte
	var ok bool

	if pkg == "" {
		log.Fatalf("pkg is empty")
	}

	if j, ok = jmap[k]; ok != true {
		log.Fatalf("%v not found in jmap", k)
	}

	var u *user.User

	u, err := user.Current()

	if err != nil {
		log.Fatalf("Unable to identify current user (%v)", err.Error())
	}

	tb := path.Join(u.HomeDir, "tmpgis") // test base

	td := path.Join(tb, pkg) // test dir

	if err = os.MkdirAll(td, 0777); err != nil {
		log.Fatalf("Unable to create %v", td)
	}

	tg := path.Join(td, "gisrc.json") // test gisrc

	j = bytes.Replace(j, []byte("SETPATH"), []byte(td), -1)

	if err = ioutil.WriteFile(tg, j, 0777); err != nil {
		log.Fatalf("Unable to write to %v (%v)", tg, err.Error())
	}

	return tg, func() { os.RemoveAll(tb) }
}

// Result holds the expected values for a zone.
type Result struct {
	User, Remote, Workspace string
	Repos                   []string
}

// Results is a collection of Result structs.
type Results []Result

// Resulter returns expected results for testing.
func Resulter(k string) Results {

	if _, ok := rmap[k]; ok != true {
		log.Fatalf("%v not found in rmap", k)
	}

	return rmap[k]
}
