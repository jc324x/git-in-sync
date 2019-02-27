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

// model is an example repository
type model struct {
	name   string // gis-Ahead
	remote string // jychri/gis-Ahead
	dir    string // /Users/jychri/tmpgis/atp/staging
}

// set m.dir to tdir + m.name or mdir + m.name.
// m.dir is used with cmd.Dir.
func (m *model) set(dir string) {
	m.dir = path.Join(dir, m.name) // set m.dir
}

// make directory fresh
func (m *model) mkdirf() {
	os.RemoveAll(m.dir)      // verify rm -rf m.dir
	os.MkdirAll(m.dir, 0766) // mkdir m.dir
}

// git init
func (m *model) init() {
	cmd := exec.Command("git", "init") // git init
	cmd.Dir = m.dir                    // set dir
	cmd.Run()                          // run
}

// hub delete -y m.remote followed by hub create
func (m *model) hub() {
	cmd := exec.Command("hub", "delete", "-y", m.remote) // verify hub delete
	cmd.Dir = m.dir                                      //
	cmd.Run()                                            //
	cmd = exec.Command("hub", "create")                  // hub create
	cmd.Dir = m.dir                                      //
	cmd.Run()                                            //
}

// create file with lorem ipsum
func (m *model) create(name string) {
	lorem := "Lorem ipsum dolor sit amet, consectetur adipiscing elit."
	data := []byte(lorem)
	name = strings.Join([]string{name, ".md"}, "")
	file := path.Join(m.dir, name)

	if err := ioutil.WriteFile(file, data, 0777); err != nil {
		log.Fatal(err)
	}
}

// git add *
func (m *model) add() {
	cmd := exec.Command("git", "add", "*") // git add
	cmd.Dir = m.dir                        // set dir
	cmd.Run()                              // run
}

// git commit -m $message
func (m *model) commit(message string) {
	cmd := exec.Command("git", "commit", "-m", message) // git commit
	cmd.Dir = m.dir                                     // set dir
	cmd.Run()                                           // run
}

// git push -u origin master
func (m *model) push() {
	cmd := exec.Command("git", "push", "-u", "origin", "master") // push -u origin master
	cmd.Dir = m.dir                                              // set dir
	cmd.Run()                                                    // run
}

// set m.dir directly
func (m *model) direct(dir string) {
	m.dir = dir
}

// hub clone $m.remote
func (m *model) clone() {
	cmd := exec.Command("hub", "clone", m.remote) // hub clone
	cmd.Dir = m.dir                               // set dir
	cmd.Run()                                     // run
}

// set model behind
func (m *model) behind() {

	if !strings.Contains(m.name, "Behind") {
		return
	}

	m.create("BEHIND")
	m.add()
	m.commit("'BEHIND' commit")
	m.push()
}

// set model ahead
func (m *model) ahead() {

	if !strings.Contains(m.name, "Ahead") {
		return
	}

	m.create("AHEAD")
	m.add()
	m.commit("'AHEAD' commit")
}

// make model untracked
func (m *model) untracked() {

	if !strings.Contains(m.name, "Untracked") {
		return
	}

	m.create("UNTRACKED")
}

// overwrite README.md with fox text
func (m *model) overwrite() {
	fox := "The sly brown fox jumped over the lazy dog."
	data := []byte(fox)
	file := path.Join(m.dir, "README.md")

	if err := ioutil.WriteFile(file, data, 0777); err != nil {
		log.Fatal(err)
	}
}

// make model dirty
func (m *model) dirty() {

	if !strings.Contains(m.name, "Dirty") {
		return
	}

	m.overwrite()
}

// hub delete -y $m.remote
func (m *model) remove() {
	cmd := exec.Command("hub", "delete", "-y", m.remote) // hub delete
	cmd.Run()                                            // run
}

type models []*model

// create model repos, configured to their names
func (ms models) startup(mdir string, tdir string) {
	var wg sync.WaitGroup
	for i := range ms {
		wg.Add(1)
		go func(m *model) {
			defer wg.Done()
			m.set(mdir)                // set models dir
			m.mkdirf()                 // make directory fresh
			m.init()                   // git init
			m.hub()                    // verify fresh repo on GitHub
			m.create("README")         // create readme
			m.add()                    // add *
			m.commit("Initial commit") // commit -m "Initial commit"
			m.push()                   // push -u origin master
			m.direct(tdir)             // set tmpgis dir
			m.clone()                  // clone to tmpgis dir
			m.set(mdir)                // set models dir
			m.behind()                 // set *Behind* models behind origin master
			m.set(tdir)                // switch to tmpgis directory
			m.ahead()                  // set *Ahead* models behind origin master
			m.untracked()              // make *Untracked* models untracked
			m.dirty()                  // make *Dirty* models dirty
		}(ms[i])
	}
	wg.Wait()
}

// remove all remotes
func (ms models) cleanup() {
	var wg sync.WaitGroup
	for i := range ms {
		wg.Add(1)
		go func(m *model) {
			defer wg.Done()
			m.remove() // remove GitHub remote
		}(ms[i])
	}
	wg.Wait()
}

// create models from set options
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
	}
	return ms
}

func user() string {

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

// create subdirectories
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

// gisrc writes a gisrc to file, with data from jmap matching key k.
func gisrc(dir string, k string) string {
	var json []byte
	var ok bool

	if json, ok = Jmap[k]; ok != true {
		log.Fatalf("%v not found in jmap", k)
	}

	gisrc := path.Join(dir, "gisrc.json")

	json = bytes.Replace(json, []byte("SETPATH"), []byte(dir), -1)

	if err := ioutil.WriteFile(gisrc, json, 0777); err != nil {
		log.Fatalf("Unable to write to %v (%v)", gisrc, err.Error())
	}

	return gisrc
}

// Public

// Jmap maps a string to JSON data as a byte slice.
var Jmap = map[string][]byte{
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

// Rmap maps a string to Results.
var Rmap = map[string]Results{
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

// Result is the expected value for a zone.
type Result struct {
	User, Remote, Workspace string
	Repos                   []string
}

// Results collects Result structs.
type Results []Result

// Resulter returns Results from map Rmap.
func Resulter(key string) Results {

	if _, ok := Rmap[key]; ok != true {
		log.Fatalf("%v not found in rmap", key)
	}

	return Rmap[key]
}

// Setup creates a test environment at ~/tmpgis/$scope/, returning
// the absolute path of ~/tmpgis/$scope/gisrc.json and a cleanup function.
func Setup(scope string, key string) (string, func()) {
	base, dir := paths(scope) // base and directory paths
	gisrc := gisrc(dir, key)  // write temporary gisrc

	return gisrc, func() {
		os.RemoveAll(base) // rm -rf base
	}
}

// Hub creates a test environment at ~/tmpgis/$scope, returning
// the path of ~/tmpgis/$scope/gisrc.json and a cleanup function.
// Unlike Setup, Hub creates matching remotes on GitHub and sets repos
// ahead, behind, dirty, untracked or complete depending according to their names.
func Hub(scope string, key string) (string, func()) {
	user := user()             // read user ~/.config/hub
	base, dir := paths(scope)  // base and directory paths (return base for cleanup...)
	gisrc := gisrc(dir, key)   // create gisrc, return it's path
	mdir, tdir := subdirs(dir) // create subdirectories models and tmp
	ms := modeler(user)        // create basic models, collect as models
	ms.startup(mdir, tdir)     // async startup

	return gisrc, func() {
		os.RemoveAll(base) // rm -rf base
		ms.cleanup()       // async remove all remotes
	}
}

// Direct verifies ~/.gisrc.json, but does not validate its content.
// If no ~/.gisrc.json, Direct creates one and returns its absolute
// path with a cleanup function. If present, Direct returns its
// absolute path with a nil cleanup function.
func Direct(scope string, key string) (string, func()) {

	var json []byte
	var ok bool

	if scope == "" {
		log.Fatalf("pkg is empty")
	}

	if json, ok = Jmap[key]; ok != true {
		log.Fatalf("%v not found in jmap", key)
	}

	tg := tilde.Abs("~/.gisrc.json") // test gisrc
	tb := tilde.Abs("~/tmpgis")      // test base
	td := path.Join(tb, scope)       // test directory

	if _, err := os.Stat(tg); err == nil {
		return "", func() {} // .gisrc.json exists, return
	}

	json = bytes.Replace(json, []byte("SETPATH"), []byte(td), -1)

	if err := ioutil.WriteFile(tg, json, 0777); err != nil {
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
