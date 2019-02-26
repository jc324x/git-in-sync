// Package atp manages test environments for git-in-sync packages.
package atp

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"

	"github.com/jychri/git-in-sync/pkg/brf"
	"github.com/jychri/git-in-sync/pkg/tilde"
)

// private
const base = "~/tmpgis"

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
	"tmpgis": []byte(`
	{
		"bundles": [{
			"path": "SETPATH",
			"zones": [{
					"user": "jychri",
					"remote": "github",
					"workspace": "tmpgis",
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
	"tmpgis": {
		{"jychri", "github", "tmpgis", []string{
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

type model struct {
	name   string // gis-Ahead
	remote string // jychri/gis-Ahead
	dir    string // /Users/jychri/tmpgis/atp/staging
}

func (m *model) hub() {
	cmd := exec.Command("hub", "delete", "-y", m.remote) // verify hub delete
	cmd.Run()                                            //
	cmd = exec.Command("hub", "create")                  // hub create
	cmd.Run()                                            //
}

func (m *model) set(dir string) {
	m.dir = path.Join(dir, m.name) // set m.dir
}

func (m *model) sdir() {
	os.RemoveAll(m.dir)      // verify rm -rf m.dir
	os.MkdirAll(m.dir, 0766) // mkdir m.dir
}

func (m *model) init() {
	cmd := exec.Command("git", "init") // git init
	cmd.Dir = m.dir                    // set dir
	cmd.Run()                          // run
}

// create name.md at m.dir with lorem ipsum
func (m *model) create(name string) {

	lorem := `Lorem ipsum dolor sit amet, consectetur adipiscing elit. \n 
	Morbi ac vulputate mi, sit amet euismod nibh. Donec at interdum sapien, 
	non pretium tortor. Duis et dapibus eros. Sed tempus non dui vel maximus. 
	\n Vivamus faucibus tellus in scelerisque ultrices. 
	Duis ac libero a leo gravida convallis. Aliquam viverra lacinia arcu, 
	ac metus pharetra sit amet. `

	data := []byte(lorem)

	name = strings.Join([]string{name, ".md"}, "")
	file := path.Join(m.dir, name)

	if err := ioutil.WriteFile(file, data, 0777); err != nil {
		log.Fatal(err)
	}
}

// ovewrite README.md with fox stuff
func (m *model) overwrite() {
	fox := "The sly brown fox jumped over the lazy dog."

	data := []byte(fox)

	file := path.Join(m.dir, "README.md")
	ioutil.WriteFile(file, data, 0777)
}

func (m *model) add() {
	cmd := exec.Command("git", "add", "*") // git add
	cmd.Dir = m.dir                        // set dir
	cmd.Run()                              // run
}

func (m *model) commit(message string) {
	cmd := exec.Command("git", "commit", "-m", message) // git commit
	cmd.Dir = m.dir                                     // set dir
	cmd.Run()                                           // run
}

func (m *model) push() {
	cmd := exec.Command("git", "push", "-u", "origin", "master") // push -u origin master
	cmd.Dir = m.dir                                              // set dir
	cmd.Run()                                                    // run
}

func (m *model) clone() {
	cmd := exec.Command("hub", "clone", m.remote) // hub clone
	cmd.Dir = m.dir                               // set dir
	cmd.Run()                                     // run
}

func (m *model) behind() {

	if !strings.Contains(m.name, "Behind") {
		return
	}

	m.create("BEHIND")
	m.add()
	m.commit("'BEHIND' commit")
	m.push()
}

func (m *model) ahead() {

	if !strings.Contains(m.name, "Ahead") {
		return
	}

	m.create("AHEAD")
	m.add()
	m.commit("'AHEAD' commit")
}

func (m *model) dirty() {

	if !strings.Contains(m.name, "Dirty") {
		return
	}

	m.overwrite()
}

func (m *model) untracked() {

	if !strings.Contains(m.name, "Untracked") {
		return
	}

	m.create("UNTRACKED")
}

func (m *model) remove() {
	cmd := exec.Command("hub", "delete", "-y", m.remote) // hub delete
	cmd.Run()                                            // run
}

// Models ...
type models []*model

func (ms models) startup(mdir string, tdir string) {
	var wg sync.WaitGroup
	for i := range ms {
		wg.Add(1)
		go func(m *model) {
			defer wg.Done()
			m.hub()                    // verify fresh repo on GitHub
			m.set(mdir)                // switch to model directory
			m.sdir()                   // verify fresh subdirectory m.dir
			m.init()                   // git init
			m.create("README")         // touch readme
			m.add()                    // add *
			m.commit("Initial commit") // commit -m "Initial commit"
			m.push()                   // push -u origin master
			m.behind()                 // set *Behind* models behind origin master
			// m.set(tdir)                // switch to tmp directory
			// m.ahead()                  // set *Ahead* models behind origin master
			// m.dirty()                  // make *Dirty* models dirty
			// m.untracked()              // make *Untracked* models untracked
		}(ms[i])
	}
	wg.Wait()
}

func (ms models) cleanup() {
	var wg sync.WaitGroup
	for i := range ms {
		wg.Add(1)
		go func(m *model) {
			defer wg.Done()
			m.remove() // remove remote origin on GitHub
		}(ms[i])
	}
	wg.Wait()
}

func modeler(user string) (ms models) {

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

	for _, name := range tmps {
		remote := path.Join(user, name)
		m := new(model)
		m.name = name
		m.remote = remote
		ms = append(ms, m)
		// log.Println(m.name)
	}
	return ms
}

func user() string {

	// also check to be sure that hub is installed...

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

	base := tilde.Abs(base)
	dir := path.Join(base, pkg)

	os.RemoveAll(dir)

	if err := os.MkdirAll(dir, 0777); err != nil {
		log.Fatalf("Unable to create %v", dir)
	}

	return base, dir
}

// models, tmps := subdirs(dir) // create sub dirs
func subdirs(dir string) (mdir string, tdir string) {
	mdir = path.Join(dir, "models")
	tdir = path.Join(dir, "tmpgis")

	if err := os.MkdirAll(mdir, 0777); err != nil {
		log.Fatalf("Unable to create %v", mdir)
	}

	if err := os.MkdirAll(tdir, 0777); err != nil {
		log.Fatalf("Unable to create %v", tdir)
	}

	return mdir, tdir
}

func fox() []byte {
	fox := "The sly brown fox jumped over the lazy dog."

	return []byte(fox)
}

// gisrcer writes a gisrc to file, data from jmap matching key k.
// add second return of function (can be nil) that will cleanup created gisrc
// if it didn't exist already
func gisrc(dir string, k string) string {
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

// Public

// Setup creates a test environment at ~/tmpgis/$pkg/.
// ~/tmpgis/$pkg/ and ~/tmpgis/$pkg/gisrc.json are created,
// key k is matched to Jmap, returning j ([]byte) if valid,
// $td replaces 'SETPATH' in j, which is written to gisrc.json.
// Setup returns the absolute path of ~/tmpgis/$pkg/gisrc.json
// and a cleanup function that removes ~/tmpgis/$pkg/.
// Note: Look at spec doc for os.MkdirAll and pull in.
func Setup(scope string, k string) (string, func()) {
	base, dir := paths(scope) // base and directory paths
	gisrc := gisrc(dir, k)    // write temporary gisrc
	return gisrc, func() { os.RemoveAll(base) }
}

// Hub ...
func Hub(scope string, k string) (string, func()) {
	user := user()             // read user ~/.config/hub
	base, dir := paths(scope)  // base and directory paths (return base for cleanup...)
	gisrc := gisrc(dir, k)     // create gisrc, return it's path
	mdir, tdir := subdirs(dir) // create subdirectories models and tmp
	ms := modeler(user)        // create basic models, collect as models
	ms.startup(mdir, tdir)     // async startup

	return gisrc, func() {
		os.RemoveAll(base) // rm -rf base
		ms.cleanup()       // async remove all remotes
	}
}

// Direct verifies the existence of a gisrc.json at ~/.gisrc.json;
// the files contents are not validated. If the path is empty,
// Direct creates a ~/.gisrc.json and returns its absolute path
// with a clean up function that will remove it. If ~/.gisrc.json
// is already present, the absolute path of ~/.gisrc.json
// is returned with a nil cleanup function.
func Direct(pkg string, k string) (string, func()) {

	var j []byte
	var ok bool

	if pkg == "" {
		log.Fatalf("pkg is empty")
	}

	if j, ok = jmap[k]; ok != true {
		log.Fatalf("%v not found in jmap", k)
	}

	tg := tilde.Abs("~/.gisrc.json") // test gisrc
	tb := tilde.Abs("~/tmpgis")      // test base
	td := path.Join(tb, pkg)         // test directory

	if _, err := os.Stat(tg); err == nil {
		return "", func() {} // .gisrc.json exists; get out
	}

	j = bytes.Replace(j, []byte("SETPATH"), []byte(td), -1)

	if err := ioutil.WriteFile(tg, j, 0777); err != nil {
		log.Fatalf("Unable to write to %v (%v)", tg, err.Error())
	}

	if err := os.MkdirAll(td, 0777); err != nil {
		log.Fatalf("Unable to create %v", td)
	}

	return tg, func() {
		os.Remove(tg)
		os.RemoveAll(tb)
	}
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
