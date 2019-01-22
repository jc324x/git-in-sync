package repos

import (
	"github.com/jychri/git-in-sync/pkg/brf"
	"github.com/jychri/git-in-sync/pkg/conf"
	"github.com/jychri/git-in-sync/pkg/repo"
	"sort"
	"sync"
)

// initRepo returns a *Repo with initial values set.

// --> Repos: Collection of Repos

type Repos []*repo.Repo

func Init(c conf.Config) (rs Repos) {

	// print
	// targetPrintln(f, "%v parsing divs|repos", e.Pager)
	// emoji.Eprint("%v parsing divs|repos", "Boat")

	// initialize Repos from Config
	for _, bl := range c.Bundles {
		for _, z := range bl.Zones {
			for _, rn := range z.Repos {
				r := repo.Init(z.Workspace, z.User, z.Remote, bl.Path, rn)
				rs = append(rs, r)
			}
		}
	}

	// timer
	// t.MarkMoment("init-repos")

	// sort
	rs.sortByPath()

	// get all divs, remove duplicates
	var dvs []string // divs

	for _, r := range rs {
		dvs = append(dvs, r.DivPath)
	}

	dvs = brf.RemoveDuplicates(dvs)

	// print
	// targetPrintln(f, "%v [%v|%v] divs|repos {%v / %v}", e.FaxMachine, len(dvs), len(rs), t.GetSplit(), t.GetTime())

	return rs
}

func initPendingRepos(rs Repos) (prs Repos) {
	for _, r := range rs {
		if r.Category == "Pending" {
			prs = append(prs, r)
		}
	}
	return prs
}

func initScheludedRepos(rs Repos) (srs Repos) {
	for _, r := range rs {
		if r.Category == "Scheduled" {
			srs = append(srs, r)
		}
	}
	return srs
}

func initSkippedRepos(rs Repos) (skrs Repos) {
	for _, r := range rs {
		if r.Category == "Skipped" {
			skrs = append(skrs, r)
		}
	}
	return skrs
}

// sort A-Z by r.Name
func (rs Repos) sortByName() {
	sort.SliceStable(rs, func(i, j int) bool { return rs[i].Name < rs[j].Name })
}

// sort A-Z by r.DivPath, then r.Name
func (rs Repos) sortByPath() {
	sort.SliceStable(rs, func(i, j int) bool { return rs[i].Name < rs[j].Name })
	sort.SliceStable(rs, func(i, j int) bool { return rs[i].DivPath < rs[j].DivPath })
}

func (rs Repos) verifyCloned() {
	var pc []string // pending clone

	for _, r := range rs {
		// r.gitCheckPending(e, f)

		if r.PendingClone == true {
			pc = append(pc, r.Name)
		}
	}

	// return if there are no pending repos

	if len(pc) <= 0 {
		return
	}

	// if there are pending repos
	// targetPrintln(f, "%v cloning [%v]", e.Sheep, len(pc))

	// verify each repo (async)
	var wg sync.WaitGroup
	for i := range rs {
		wg.Add(1)
		go func(r *repo.Repo) {
			defer wg.Done()
			// r.gitClone(e, f)
		}(rs[i])
	}
	wg.Wait()

	var cr []string // cloned repos

	for _, r := range rs {
		if r.Cloned == true {
			cr = append(cr, r.Name)
		}
	}

	// timer
	// t.MarkMoment("verify-repos")

	// summary
	// var b bytes.Buffer

	// b.WriteString(e.Truck)
	// b.WriteString(" [")
	// b.WriteString(strconv.Itoa(len(cr)))
	// b.WriteString("/")
	// b.WriteString(strconv.Itoa(len(pc)))
	// b.WriteString("] cloned")

	// tr := time.Millisecond // truncate
	// b.WriteString(" {")
	// b.WriteString(t.GetSplit().Truncate(tr).String())
	// b.WriteString(" / ")
	// b.WriteString(t.GetTime().Truncate(tr).String())
	// b.WriteString("}")

	// targetPrintln(f, b.String())
}
