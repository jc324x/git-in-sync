package repos

import (
	"log"
	// "os"
	"path"
	"path/filepath"
	"testing"

	// "github.com/jychri/git-in-sync/pkg/brf"
	"github.com/jychri/git-in-sync/pkg/conf"
	// "github.com/jychri/git-in-sync/pkg/emoji"
	// "github.com/jychri/git-in-sync/pkg/fchk"
	"github.com/jychri/git-in-sync/pkg/flags"
	"github.com/jychri/git-in-sync/pkg/timer"
)

// func TestInit(t *testing.T) {

// 	var abs string
// 	var err error

// 	if abs, err = filepath.Abs(""); err != nil {
// 		t.Errorf("TestInit: err = %v\n", err.Error())
// 	}

// 	p := path.Join(abs, "ex_gisrc.json")

// 	json := []byte(`
// 		{
// 			"bundles": [{
// 				"path": "~/testing_gis",
// 				"zones": [{
// 						"user": "hendricius",
// 						"remote": "github",
// 						"workspace": "recipes",
// 						"repositories": [
// 							"pizza-dough",
// 							"the-bread-code"
// 						]
// 					},
// 					{
// 						"user": "cocktails-for-programmers",
// 						"remote": "github",
// 						"workspace": "recipes",
// 						"repositories": [
// 							"cocktails-for-programmers"
// 						]
// 					},
// 					{
// 						"user": "rochacbruno",
// 						"remote": "github",
// 						"workspace": "recipes",
// 						"repositories": [
// 							"vegan_recipes"
// 						]
// 					},
// 					{
// 						"user": "niw",
// 						"remote": "github",
// 						"workspace": "recipes",
// 						"repositories": [
// 							"ramen"
// 						]
// 					}
// 				]
// 			}]
// 		}
// `)

// 	if err = ioutil.WriteFile(p, json, 0644); err != nil {
// 		t.Errorf("Init (%v)", err.Error())
// 	}

// 	f := flags.Flags{Mode: "verify", Config: p}

// 	c := Init(f)

func setup(c conf.Config, f flags.Flags, t *timer.Timer) {

	var abs string
	var err error

	if abs, err = filepath.Abs(""); err != nil {
		log.Fatalf("setup (%v)", err.Error())
	}

	p := path.Join(abs, "test")

	json := []byte(`
		{
			"bundles": [{
				"path": "~/git-in-sync-tests",
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
`)

	f = flags.Flags{Mode: "verify", Config: p}
	f.Config = p

}

func TestInit(t *testing.T) {

}

// func TestVerifyDivs(t *testing.T) {

// }
