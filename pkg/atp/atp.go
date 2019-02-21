// Package atp manages test environments for git-in-sync packages.
package atp

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"sync"

	"github.com/jychri/git-in-sync/pkg/brf"
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
`),
	"tmp": []byte(`
	{
		"bundles": [{
			"path": "SETPATH",
			"zones": [{
					"user": "jychri",
					"remote": "github",
					"workspace": "tmp",
					"repositories": [
						"gis-Ahead",
						"gis-Behind",
						"gis-Dirty",
						"gis-DirtyUntracked",
						"gis-DirtyAhead",
						"gis-DirtyBehind",
						"gis-Untracked",
						"gis-UntrackedAhead",
						"gis-UntrackedBehind",
						"gis-Complete"
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
	"tmp": {
		{"jychri", "github", "tmp", []string{
			"gis-Ahead",
			"gis-Behind",
			"gis-Dirty",
			"gis-DirtyUntracked",
			"gis-DirtyAhead",
			"gis-DirtyBehind",
			"gis-Untracked",
			"gis-UntrackedAhead",
			"gis-UntrackedBehind",
			"gis-Complete",
		}}},
}

// tmp repo names. Hub creates these repos on disk and remotes on GitHub...
// tmps are matched to Results and JSON config in rmap and jmap...
var tmps = []string{
	"gis-Ahead",
	"gis-Behind",
	"gis-Dirty",
	"gis-DirtyUntracked",
	"gis-DirtyAhead",
	"gis-DirtyBehind",
	"gis-Untracked",
	"gis-UntrackedAhead",
	"gis-UntrackedBehind",
	"gis-Complete",
}

// Tmp is holds tmp repository details.
type Tmp struct {
	Category string // main or exp? along those lines
	Path     string // absolute path on disk
	Remote   string // $user/$repo
}

func hubconfig() string {

	path := tilde.Abs("~/.config/hub")
	file, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var user string // user set in ~/.config/hub
	var token bool  // token present in ~/.config/hub

	for scanner.Scan() {
		l := scanner.Text()

		if match := brf.MatchLine(l, "- user:"); match != "" {
			user = match
		}

		if match := brf.MatchLine(l, "oauth_token:"); match != "" {
			token = true
		}
	}

	if user == "" || token == false {
		log.Fatalf("Error in ~/.config.hub")
	}

	return user
}

func paths(pkg string) (string, string) {
	if pkg == "" {
		log.Fatalf("pkg is empty")
	}

	base := tilde.Abs("~/tmpgis")
	dir := path.Join(base, pkg)

	if err := os.MkdirAll(dir, 0777); err != nil {
		log.Fatalf("Unable to create %v", dir)
	}

	return base, dir
}

// write writes a gisrc to file
func write(dir string, k string) string {
	var j []byte
	var ok bool

	if j, ok = jmap[k]; ok != true {
		log.Fatalf("%v not found in jmap", k)
	}

	gisrc := path.Join(dir, "gisrc.json")

	j = bytes.Replace(j, []byte("SETPATH"), []byte(dir), -1)

	if err := ioutil.WriteFile(gisrc, j, 0777); err != nil {
		log.Fatalf("Unable to write to %v (%v)", gisrc, err.Error())
	}

	return gisrc
}

// startup creates a new temporary repo in the local test directory
// 'tmp'. The repo is initialized with a blank README.md file, a first
// commit with a message of 'Initial commit' and a upstream master branch
// available at github.com/$user/$tmp. startup is called by Hub, which returns
// a cleanup function that deletes the remote branch and deletes all temp
// repos and directories.
func startup(dir string, user string, tmp string) string {

	local := path.Join(dir, "tmp", tmp)

	os.RemoveAll(local)      // remove
	os.MkdirAll(local, 0777) // create

	// git init
	cmd := exec.Command("git", "init")
	cmd.Dir = local
	cmd.Run()

	// hub create
	cmd = exec.Command("hub", "create")
	cmd.Dir = local
	cmd.Run()

	// touch README.md
	cmd = exec.Command("touch", "README.md")
	cmd.Dir = local
	cmd.Run()

	// git add *
	cmd = exec.Command("git", "add", "*")
	cmd.Dir = local
	cmd.Run()

	// git commit -m "Initial commit"
	cmd = exec.Command("git", "commit", "-m", "Initial commit")
	cmd.Dir = local
	cmd.Run()

	// git commit -- set-upstream origin master
	cmd = exec.Command("git", "push", "--set-upstream", "origin", "master")
	cmd.Dir = local
	cmd.Run()

	return path.Join(user, tmp)
}

// create all remotes on GitHub
func create(tmps []string, dir string, user string) (remotes []string) {
	var wg sync.WaitGroup

	for i := range tmps {
		wg.Add(1)
		go func(tmp string) {
			defer wg.Done()
			remote := startup(dir, user, tmp)
			remotes = append(remotes, remote)
		}(tmps[i])
	}
	wg.Wait()

	return remotes
}

// remove (delete) all remotes on GitHub
func remove(remotes []string) {
	for _, remote := range remotes {
		cmd := exec.Command("hub", "delete", "-y", remote)
		cmd.Run()
	}
}

// commit creates an additional commit to a repo at path
func commit() {

}

func dirty() {

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
	base, dir := paths(pkg)
	gisrc := write(dir, k)
	return gisrc, func() { os.RemoveAll(base) }
}

// Base returns the base testing directory
func Base(pkg string) string {
	base, _ := paths(pkg)
	return base
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

	return tg, func() { os.Remove(tg) }
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

// Hub uses hub to do stuff...
// this needs to return a string path to a .gisrc just like Setup
// should be ~/tmpgis/tmp/
func Hub(pkg string, k string) (string, string, func()) {

	base, dir := paths(pkg)            // base and directory paths
	gisrc := write(dir, k)             // create the main gisrc
	user := hubconfig()                // read the user from ~/.config/hub
	remotes := create(tmps, dir, user) // creates repos and remotes on GitHub
	exp := fmt.Sprintf("%v-exp", pkg)  // secondary exp
	_, dir = paths(exp)                // exp path
	egisrc := write(dir, k)            // create secondary gisrc in exp

	// cleanup function
	return gisrc, egisrc, func() {
		os.RemoveAll(base)
		remove(remotes)
	}
}

// Conditions sets Setup
func Conditions() {
	// @ exp path, get all directories
	// if name has Behind, drop in a commit and push to remote

	// @ main path, get all directories
	// if
}
