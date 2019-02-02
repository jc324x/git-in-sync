package repos

import (
	"testing"

	"github.com/jychri/git-in-sync/pkg/atp"
	"github.com/jychri/git-in-sync/pkg/conf"
	"github.com/jychri/git-in-sync/pkg/flags"
	"github.com/jychri/git-in-sync/pkg/timer"
)

func TestVerify(t *testing.T) {
	p, _ := atp.Setup("testing", "recipes")
	ti := timer.Init()
	f := flags.Testing(p)
	c := conf.Init(f)
	rs := Init(c, f, ti)

	// defer cleanup()

	rs.VerifyWorkspaces(f, ti)
}

func TestAgain(t *testing.T) {
	_, cleanup := atp.Setup("testing", "recipes")
	// ti := timer.Init()
	// // f := flags.Testing(p)
	// c := conf.Init(f)
	// rs := Init(c, f, ti)

	defer cleanup()
}
