package repo

import (
	"testing"
)

func TestInit(t *testing.T) {

	r := Init("main", "jychri", "github", "~/tmpgis", "git-in-sync")

	if r.BundlePath != "~/tmpgis" {
		t.Errorf("Init: %v != '~/tmpgis'", r.BundlePath)
	}

	if r.Workspace != "main" {
		t.Errorf("Init: %v != 'main'", r.Workspace)
	}

	if r.User != "jychri" {
		t.Errorf("Init: %v != 'jychri'", r.User)
	}

	if r.Remote != "github" {
		t.Errorf("Init: %v != 'github'", r.Remote)
	}

	if r.Name != "git-in-sync" {
		t.Errorf("Init: %v != 'git-in-sync'", r.Name)
	}

	// if r.WorkspacePath != "" {

	// }

	// if r.RepoPath != "" {

	// }

	// if r.GitPath != "" {

	// }

	// if r.GitDir != "" {

	// }

	// if r.WorkTree != "" {

	// }

	if r.URL != "https://github.com/jychri/git-in-sync" {
		t.Errorf("Init: %v != 'https://github.com/jychri/git-in-sync'", r.URL)
	}

}

// WorkspacePath    string   // "/Users/jychri/dev/go-lang/"
// RepoPath         string   // "/Users/jychri/dev/go-lang/git-in-sync"
// GitPath          string   // "/Users/jychri/dev/go-lang/git-in-sync/.git"
// GitDir           string   // "--git-dir=/Users/jychri/dev/go-lang/git-in-sync/.git"
// WorkTree         string   // "--work-tree=/Users/jychri/dev/go-lang/git-in-sync"
// URL              string   // "https://github.com/jychri/git-in-sync"
