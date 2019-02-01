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
	// "github.com/jychri/git-in-sync/pkg/test"
	"github.com/jychri/git-in-sync/pkg/timer"
)

func setup(c conf.Config, f flags.Flags, t *timer.Timer) {

	var abs string
	var err error

	if abs, err = filepath.Abs(""); err != nil {
		log.Fatalf("setup (%v)", err.Error())
	}

	p := path.Join(abs, "test")

	f = flags.Flags{Mode: "verify", Config: p}
	f.Config = p

}

func TestInit(t *testing.T) {

}

// func TestVerifyDivs(t *testing.T) {

// }
