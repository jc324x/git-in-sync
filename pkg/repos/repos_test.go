package repos

import (
	// "log"
	"testing"

	"github.com/jychri/git-in-sync/pkg/atp"
	"github.com/jychri/git-in-sync/pkg/conf"
	"github.com/jychri/git-in-sync/pkg/flags"
	"github.com/jychri/git-in-sync/pkg/timer"
)

func TestVerify(t *testing.T) {
	p, cleanup := atp.Setup("testing", "recipes")
	ti := timer.Init()
	f := flags.Testing(p)
	c := conf.Init(f)
	rs := Init(c, f, ti)

	// twsp := rs.workspaces()
	// log.Println(len(twsp))

	// tws := rs.workspacePaths()
	// log.Println(len(tws))

	defer cleanup()

	rs.VerifyWorkspaces(f, ti)
}
