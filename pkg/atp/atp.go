// Package atp manages test environments for git-in-sync packages.
package atp

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"sync"

	"github.com/jychri/git-in-sync/pkg/tilde"
)

// private

// JSON map
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
							"pizza-dough"
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

// Result map
var rmap = map[string]Results{
	"recipes": {
		{"hendricius", "github", "recipes", []string{"pizza-dough"}},
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
		{"jychri", "github", "tmp", []string{
			"tmp1",
			"tmp2",
			"tmp3",
			"tmp4",
			"tmp5",
		}}},
}

// test repos
var trs = []string{
	"tmpgis0",
	"tmpgis1",
	"tmpgis2",
	"tmpgis3",
	"tmpgis4",
	"tmpgis5",
	"tmpgis6",
	"tmpgis7",
	"tmpgis8",
	"tmpgis9",
}

// Public

// Setup creates a test environment at ~/tmpgis/$pkg/.
// ~/tmpgis/$pkg/ and ~/tmpgis/$pkg/gisrc.json are created,
// key k is matched to Jmap, returning j ([]byte) if valid,
// $td replaces 'SETPATH' in j, which is written to gisrc.json.
// Setup returns the absolute path of ~/tmpgis/$pkg/gisrc.json
// and a cleanup function that removes ~/tmpgis/$pkg/.
// Note: Look at spec doc for os.MkdirAll and pull in.
func Setup(pkg string, k string) (string, func()) {

	var j []byte
	var ok bool

	if pkg == "" {
		log.Fatalf("pkg is empty")
	}

	if j, ok = jmap[k]; ok != true {
		log.Fatalf("%v not found in jmap", k)
	}

	tb := tilde.Abs("~/tmpgis") // test base
	td := path.Join(tb, pkg)    // test dir

	if err := os.MkdirAll(td, 0777); err != nil {
		log.Fatalf("Unable to create %v", td)
	}

	tg := path.Join(td, "gisrc.json")                       // test gisrc
	j = bytes.Replace(j, []byte("SETPATH"), []byte(td), -1) // SETPATH set

	if err := ioutil.WriteFile(tg, j, 0777); err != nil {
		log.Fatalf("Unable to write to %v (%v)", tg, err.Error())
	}

	return tg, func() { os.RemoveAll(tb) }
}

// Directory returns that path of testing environment ~/tmpgis/$pkg/.
func Directory(pkg string) string {

	if pkg == "" {
		log.Fatalf("pkg is empty")
	}

	tb := tilde.Abs("~/tmpgis")

	return path.Join(tb, pkg)
}

// Direct verifies the existence of a gisrc.json at ~/.gisrc.json;
// the files contents are not validated. If the path is empty,
// Direct creates a ~/.gisrc.json and returns its absolute path
// with a clean up function that removes it. If ~/.gisrc.json
// is present, the absolute path of ~/.gisrc.json
// is returned with a mute cleanup function.
func Direct(pkg string, k string) (string, func()) {

	var j []byte
	var ok bool

	if pkg == "" {
		log.Fatalf("pkg is empty")
	}

	if j, ok = jmap[k]; ok != true {
		log.Fatalf("%v not found in jmap", k)
	}

	tg := tilde.Abs("~/.gisrc.json")

	if _, err := os.Stat(tg); err == nil {
		return "", func() {} // .gisrc.json exists; get out
	}

	tb := tilde.Abs("~/tmpgis")

	td := path.Join(tb, pkg)

	j = bytes.Replace(j, []byte("SETPATH"), []byte(td), -1)

	if err := ioutil.WriteFile(tg, j, 0777); err != nil {
		log.Fatalf("Unable to write to %v (%v)", tg, err.Error())
	}

	if err := os.MkdirAll(td, 0777); err != nil {
		log.Fatalf("Unable to create %v", td)
	}

	return tg, func() { log.Printf("removing %v\n", tg); os.Remove(tg) }
}

// Result holds the expected values for a zone.
type Result struct {
	User, Remote, Workspace string
	Repos                   []string
}

// Results collects Result structs.
type Results []Result

// Resulter returns expected results as Results
// from map rmap. Given an unrecognized key,
// execution is stopped with log.Fatalf().
func Resulter(k string) Results {

	if _, ok := rmap[k]; ok != true {
		log.Fatalf("%v not found in rmap", k)
	}

	return rmap[k]
}

// Hub uses GitHub's binary to create repos
// with remotes.
func Hub(pkg string) (string, func()) {

	var trps []string

	if pkg == "" {
		log.Fatalf("pkg is empty")
	}

	tb := tilde.Abs("~/tmpgis") // test base
	td := path.Join(tb, pkg)    // test dir

	if err := os.MkdirAll(td, 0777); err != nil {
		log.Fatalf("Unable to create %v", td)
	}

	var wg sync.WaitGroup
	for i := range trs {
		wg.Add(1)
		go func(tr string) {
			defer wg.Done()

			tp := path.Join(td, tr) // test path

			// mkdir
			if _, err := os.Stat(tp); os.IsNotExist(err) {
				log.Printf("creating %v\n", tp)
				os.MkdirAll(tp, 0777)
			}

			// git init
			cmd := exec.Command("git", "init")
			log.Printf("init %v\n", tp)
			cmd.Dir = tp
			cmd.Run()

			// hub create
			cmd = exec.Command("hub", "create")
			log.Printf("hub create %v\n", tp)
			cmd.Dir = tp
			cmd.Run()

			// touch README.md
			cmd = exec.Command("touch", "README.md")
			log.Printf("touch %v\n", tp)
			cmd.Dir = tp
			cmd.Run()

			// git add *
			cmd = exec.Command("git", "add", "*")
			log.Printf("git add * %v\n", tp)
			cmd.Dir = tp
			cmd.Run()

			// git commit -m "Initial commit"
			cmd = exec.Command("git", "commit", "-m", "Initial commit")
			log.Printf("touch %v\n", tp)
			cmd.Dir = tp
			cmd.Run()

			// git commit -- set-upstream origin master
			cmd = exec.Command("git", "push", "--set-upstream", "origin", "master")
			log.Printf("push upstream %v\n", tp)
			cmd.Dir = tp
			cmd.Run()

			// add to trp for removal later
			trps = append(trps, tp)

		}(trs[i])
	}
	wg.Wait()

	return td, func() {
		os.RemoveAll(tb)

		for _, trp := range trps {
			// hub delet
			cmd := exec.Command("hub", "delete", "yes")
			log.Printf("hub delete %v\n", trp)
			cmd.Dir = trp
			cmd.Run()
		}
	}
}
