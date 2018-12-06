package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// --> Moment: a moment in time

type Moment struct {
	Name  string
	Time  time.Time
	Start time.Duration // duration since start
	Split time.Duration // duration since last moment
}

// --> Timer: tracking moments in time

type Timer struct {
	Moments []Moment
}

// initTimer initializes a *Timer with a Start moment.
func initTimer() *Timer {
	t := new(Timer)
	st := Moment{Name: "Start", Time: time.Now()} // (st)art
	t.Moments = append(t.Moments, st)
	return t
}

// markMoment marks a moment in time as a Moment and appends t.Moments.
func (t *Timer) markMoment(s string) {
	sm := t.Moments[0]                           // (s)tarting (m)oment
	lm := t.Moments[len(t.Moments)-1]            // (l)ast (m)oment
	m := Moment{Name: s, Time: time.Now()}       // name and time
	m.Start = time.Since(sm.Time).Truncate(1000) // duration since start
	m.Split = m.Start - lm.Start                 // duration since last moment
	t.Moments = append(t.Moments, m)             // append Moment
}

// getTime returns the elapsed time at the last recorded moment in t.Moments.
func (t *Timer) getTime() time.Duration {
	lm := t.Moments[len(t.Moments)-1] // (l)ast (m)oment
	return lm.Start
}

// getSplit returns the split time for the last recorded moment in t.Moments.
func (t *Timer) getSplit() time.Duration {
	lm := t.Moments[len(t.Moments)-1] // (l)ast (m)oment
	return lm.Split
}

// getMoment returns a Moment and an error value from t.Moments.
func (t *Timer) getMoment(s string) (Moment, error) {
	for _, m := range t.Moments {
		if m.Name == s {
			return m, nil
		}
	}

	var em Moment // (e)mpty (m)oment
	return em, errors.New("no moment found")
}

// --> Emoji: struct collecting emojis

type Emoji struct {
	AlarmClock           string
	Book                 string
	Books                string
	Box                  string
	BuildingConstruction string
	Bunny                string
	Checkmark            string
	Clapper              string
	Clipboard            string
	CrystalBall          string
	Desert               string
	DirectHit            string
	FaxMachine           string
	Finger               string
	Flag                 string
	FlagInHole           string
	FileCabinet          string
	Fire                 string
	Folder               string
	Glasses              string
	Hourglass            string
	Hole                 string
	Inbox                string
	Memo                 string
	Microscope           string
	Outbox               string
	Pager                string
	Parents              string
	Pen                  string
	Pig                  string
	Popcorn              string
	Rocket               string
	Run                  string
	Satellite            string
	SatelliteDish        string
	Ship                 string
	Slash                string
	Squirrel             string
	Telescope            string
	Text                 string
	ThinkingFace         string
	TimerClock           string
	Traffic              string
	Truck                string
	Turtle               string
	ThumbsUp             string
	Unicorn              string
	Warning              string
	Count                int
}

// initEmoji returns an Emoji struct with all values initialized.
func initEmoji(f Flags, t *Timer) (e Emoji) {
	e.AlarmClock = printEmoji(9200)
	e.Book = printEmoji(128214)
	e.Books = printEmoji(128218)
	e.Box = printEmoji(128230)
	e.BuildingConstruction = printEmoji(127959)
	e.Bunny = printEmoji(128048)
	e.Checkmark = printEmoji(9989)
	e.Clapper = printEmoji(127916)
	e.Clipboard = printEmoji(128203)
	e.CrystalBall = printEmoji(128302)
	e.DirectHit = printEmoji(127919)
	e.Desert = printEmoji(127964)
	e.FaxMachine = printEmoji(128224)
	e.Finger = printEmoji(128073)
	e.FileCabinet = printEmoji(128452)
	e.Flag = printEmoji(127937)
	e.FlagInHole = printEmoji(9971)
	e.Fire = printEmoji(128293)
	e.Folder = printEmoji(128193)
	e.Glasses = printEmoji(128083)
	e.Hole = printEmoji(128371)
	e.Hourglass = printEmoji(9203)
	e.Inbox = printEmoji(128229)
	e.Microscope = printEmoji(128300)
	e.Memo = printEmoji(128221)
	e.Outbox = printEmoji(128228)
	e.Pager = printEmoji(128223)
	e.Parents = printEmoji(128106)
	e.Pen = printEmoji(128394)
	e.Pig = printEmoji(128055)
	e.Popcorn = printEmoji(127871)
	e.Rocket = printEmoji(128640)
	e.Run = printEmoji(127939)
	e.Satellite = printEmoji(128752)
	e.SatelliteDish = printEmoji(128225)
	e.Slash = printEmoji(128683)
	e.Ship = printEmoji(128674)
	e.Squirrel = printEmoji(128063)
	e.Telescope = printEmoji(128301)
	e.Text = printEmoji(128172)
	e.ThumbsUp = printEmoji(128077)
	e.TimerClock = printEmoji(9202)
	e.Traffic = printEmoji(128678)
	e.Truck = printEmoji(128666)
	e.Turtle = printEmoji(128034)
	e.Unicorn = printEmoji(129412)
	e.Warning = printEmoji(128679)
	e.Count = reflect.ValueOf(e).NumField() - 1

	// timer
	t.markMoment("emoji")

	return e
}

// printEmoji returns an emoji character as a string value.
func printEmoji(n int) string {
	str := html.UnescapeString("&#" + strconv.Itoa(n) + ";")
	return str
}

// --> Flags: struct collecting flag values

type Flags struct {
	Mode    string
	Clear   bool
	Verbose bool
	Dry     bool
	Emoji   bool
	OneLine bool
	Count   int
	Summary string
}

func initFlags(e Emoji, t *Timer) (f Flags) {

	// shortcut variables
	var m string // mode
	var c bool   // clear
	var v bool   // verbose
	var d bool   // dry
	var em bool  // emoji
	var o bool   // one-line

	// summary and count
	var fc int   // flag count
	var s string // summary

	// point to shortcut variables
	flag.StringVar(&m, "m", "verify", "mode")
	flag.BoolVar(&c, "c", false, "clear")
	flag.BoolVar(&v, "v", true, "verbose")
	flag.BoolVar(&d, "d", false, "dry")
	flag.BoolVar(&em, "e", true, "emoji")
	flag.BoolVar(&o, "o", false, "one-line")
	flag.Parse()

	// collect and join (e)nabled (f)lags
	var ef []string

	// mode
	if m != "" {
		fc += 1
	}

	// ...otherwise set to 'verify'
	switch m {
	case "login", "logout", "verify":
	default:
		m = "verify"
	}
	ef = append(ef, m)

	// clear
	if c == true {
		fc += 1
		ef = append(ef, "clear")
	}

	// dry
	if d == true {
		fc += 1
		ef = append(ef, "dry")
	}

	// verbose
	if v == true {
		fc += 1
		ef = append(ef, "verbose")
	}

	// emoji
	if em == true {
		fc += 1
		ef = append(ef, "emoji")
	}

	// one-line
	if o == true {
		fc += 1
		ef = append(ef, "one-line")
	}

	// summary
	s = strings.Join(ef, ", ")

	// timer
	t.markMoment("flags")

	// set Flags
	f = Flags{m, c, v, d, em, o, fc, s}

	return f
}

// isClear returns true if f.Clear is true.
func isClear(f Flags) bool {
	if f.Clear {
		return true
	} else {
		return false
	}
}

// isVerbose returns true if f.Verbose is true.
func isVerbose(f Flags) bool {
	if f.Verbose {
		return true
	} else {
		return false
	}
}

// isDry returns true if f.Dry is true.
func isDry(f Flags) bool {
	if f.Dry {
		return true
	} else {
		return false
	}
}

// isActive returns true if f.Dry is true.
func isActive(f Flags) bool {
	if f.Dry {
		return false
	} else {
		return true
	}
}

// hasEmoji returns true if f.Emoji is true.
func hasEmoji(f Flags) bool {
	if f.Emoji {
		return true
	} else {
		return false
	}
}

// noEmoji returns true if f.Emoji is false.
func noEmoji(f Flags) bool {
	if f.Emoji {
		return false
	} else {
		return true
	}
}

// oneLine returns true if f.OneLine is true.
func oneLine(f Flags) bool {
	if f.OneLine {
		return true
	} else {
		return false
	}
}

// initPrint

func initPrint(e Emoji, f Flags, t *Timer) {
	clearScreen(f)
	targetPrint(f, "%v start", e.Clapper)

	if isDry(f) {
		targetPrint(f, "%v  dry run; no changes will be made", e.Desert)
	}

	if ft, err := t.getMoment("flags"); err == nil {
		targetPrint(f, "%v parsing flags", e.FlagInHole)
		targetPrint(f, "%v [%v] flags (%v) {%v / %v}", e.Flag, f.Count, f.Summary, ft.Split, ft.Start)
	}

	if et, err := t.getMoment("emoji"); err == nil {
		targetPrint(f, "%v initializing emoji", e.CrystalBall)
		targetPrint(f, "%v [%v] emoji {%v / %v}", e.DirectHit, e.Count, et.Split, et.Start)
	}
}

// --> Config: ~/.gisrc.json unmarshalled

type Config struct {
	Zones []struct {
		Path    string `json:"path"`
		Bundles []struct {
			User     string   `json:"user"`
			Remote   string   `json:"remote"`
			Division string   `json:"division"`
			Repos    []string `json:"repositories"`
		} `json:"bundles"`
	} `json:"zones"`
}

// initConfig returns
func initConfig(e Emoji, f Flags, t *Timer, w *Workspace) Config {

	// get the current user
	u, err := user.Current()

	if err != nil {
		log.Fatal(err)
	}

	// expand "~/" to "/Users/user"
	g := fmt.Sprintf("%v/.gisrc.json", u.HomeDir)

	// print
	targetPrint(f, "%v reading %v", e.Glasses, g)

	// read file
	r, err := ioutil.ReadFile(g)

	if err != nil {
		log.Fatalf("No file found at %v\n", g)
	}

	// unmarshall json
	c := Config{}
	err = json.Unmarshal(r, &c)

	if err != nil {
		log.Fatalf("Can't unmarshal JSON from %v\n", g)
	}

	// timer
	t.markMoment("config")

	// print
	targetPrint(f, "%v read %v {%v / %v}", e.Book, g, t.getSplit(), t.getTime())

	return c
}

// --> Repo

type Repo struct {

	// initRepo
	Div      *Div
	Name     string
	User     string
	Remote   string
	GitDir   string
	Path     string
	GitPath  string
	WorkTree string
	URL      string

	// verify
	Verified         bool
	Cloned           bool     // gitClone
	VerifiedURL      string   // gitConfigOriginURL
	RemoteUpdate     string   // gitRemoteUpdate
	Clean            bool     // gitStatusPorcelain
	LocalBranch      string   // gitLocalBranch
	LocalSHA         string   // gitLocalSHA
	MergeBaseSHA     string   // gitMergeBaseSHA
	UpstreamSHA      string   // gitUpstreamSHA
	UpstreamBranch   string   // gitRevParseUpstream
	DiffFiles        []string // gitDiffFiles
	DiffCount        int      // getDiffSummary
	DiffSummary      string   // getDiffSummary
	DiffStatus       bool     // getDiffSummary
	ShortStat        string   // gitShortstat
	ShortStatPlus    int      // getShortInts
	ShortStatMinus   int      // getShortInts
	Upstream         string   // getUpstreamStatus
	UntrackedFiles   []string // gitUntracked
	UntrackedCount   int      // getUntrackedSummary
	UntrackedSummary string   // getUntrackedSummary
	UntrackedStatus  bool     // getUntrackedSummary
	Summary          string   // getSummary
	Phase            string   // getPhase
	InfoVerified     bool     // verifyProjectInfo

	// setActions
	Status       string
	GitAction    string
	GitMessage   string
	GitConfirmed bool
}

func initRepo(d *Div, rn string, bu string, br string, bd string) *Repo {
	r := new(Repo)
	r.Div = d
	r.Name = rn
	r.User = bu   // bundle user
	r.Remote = br // bundle remote

	var b bytes.Buffer

	b.WriteString("--git-dir=")
	b.WriteString(r.Div.Path)
	b.WriteString("/")
	b.WriteString(r.Name)
	b.WriteString("/.git")
	r.GitDir = b.String()

	b.Reset()
	b.WriteString(r.Div.Path)
	b.WriteString("/")
	b.WriteString(r.Name)
	b.WriteString("/.git")
	r.GitPath = b.String()

	b.Reset()
	b.WriteString(r.Div.Path)
	b.WriteString("/")
	b.WriteString(r.Name)
	r.Path = b.String()

	b.Reset()
	b.WriteString("--work-tree=")
	b.WriteString(r.Div.Path)
	b.WriteString("/")
	b.WriteString(r.Name)
	r.WorkTree = b.String()

	b.Reset()
	switch r.Remote {
	case "github":
		b.WriteString("https://github.com/")
	case "gitlab":
		b.WriteString("https://gitlab.com/")
	}
	b.WriteString(r.User)
	b.WriteString("/")
	b.WriteString(r.Name)

	switch r.Remote {
	case "gitlab":
		b.WriteString(".git")
	}

	r.URL = b.String()

	return r
}

// repo fns here

func (r *Repo) verify(e Emoji, f Flags, w *Workspace) {

	if noDiv(r) {
		targetPrint(f, "%v %v (%v/%v)", e.Slash, r.Name, r.Div.Name, r.Remote)
		return
	}

	if canClone(f, r) {
		r.gitClone(e, f, w)
	}

	if isMissing(r) {
		targetPrint(f, "%v %v (%v/%v)", e.Slash, r.Name, r.Div.Name, r.Remote)
		return
	}

	r.gitConfigOriginURL()
	r.gitRemoteUpdate()
	r.gitStatusPorcelain()
	r.gitLocalSHA()
	r.gitLocalBranch()
	r.gitMergeBaseSHA()
	r.gitUpstreamSHA()
	r.gitRevParseUpstream()
	r.gitDiffsNameOnly()
	r.getDiffSummary()
	r.gitShortstat()
	r.getShortInts()
	r.gitUntracked()
	r.getUntrackedSummary()
	r.getUpstreamStatus()
	r.getPhase()

	if isUpToDate(r) {
		w.VerifiedRepos = append(w.VerifiedRepos, r)
		targetPrint(f, "%v %v (%v/%v)", e.Checkmark, r.Name, r.Div.Name, r.Remote)
	} else {
		// fmt.Println("DEBUG:")
		// fmt.Println(r)
		fmt.Printf("DEBUG: Phase = %v\n", r.Phase)

		targetPrint(f, "%v %v (%v/%v)", e.Warning, r.Name, r.Div.Name, r.Remote)
	}
}

func canClone(f Flags, r *Repo) bool {

	_, err := os.Stat(r.Path)

	if os.IsNotExist(err) && isActive(f) {
		return true
	} else {
		return false
	}
}

func noDiv(r *Repo) bool {
	if r.Div.PathVerified == false {
		return true
	} else {
		r.Verified = true
		return false
	}
}

func isMissing(r *Repo) bool {
	_, err := os.Stat(r.GitPath)

	if os.IsNotExist(err) {
		return true
	} else {
		r.Verified = true
		return false
	}
}

func isUpToDate(r *Repo) bool {
	if r.Verified == true {
		switch {
		case r.LocalSHA == "":
			r.InfoVerified = false
		case r.RemoteUpdate == "":
			r.InfoVerified = false
		case r.MergeBaseSHA == "":
			r.InfoVerified = false
		case r.UpstreamSHA == "":
			r.InfoVerified = false
		case r.UpstreamBranch == "":
			r.InfoVerified = false
		case r.Phase == "":
			r.InfoVerified = false
		}
		r.InfoVerified = true
	}

	if r.Verified == true && r.Phase == "Up-To-Date" {
		return true
	} else {
		return false
	}
}

func (r *Repo) gitClone(e Emoji, f Flags, w *Workspace) {
	targetPrint(f, "%v cloning %v {%v}", e.Box, r.Name, r.Div.Name)

	args := []string{"clone", r.URL, r.Path}
	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()

	w.ClonedRepos = append(w.ClonedRepos, r)
}

func (r *Repo) gitConfigOriginURL() {
	if r.Verified {
		args := []string{r.GitDir, "config", "--get", "remote.origin.url"}
		cmd := exec.Command("git", args...)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Run()
		s := out.String()
		s = strings.TrimSuffix(s, "\n")
		r.VerifiedURL = s

		if r.VerifiedURL == r.URL {
			r.Verified = true
		} else {
			r.Verified = false
		}
	}
}

func (r *Repo) gitRemoteUpdate() {
	if r.Verified {
		args := []string{r.GitDir, r.WorkTree, "fetch", "origin"}
		cmd := exec.Command("git", args...)
		var out bytes.Buffer
		var err bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &err
		cmd.Run()

		if str := out.String(); str != "" {
			r.RemoteUpdate = trim(out.String())
		}

		if str := err.String(); str != "" {
			r.RemoteUpdate = trim(err.String())
		}

	}
}

func (r *Repo) gitStatusPorcelain() {
	if r.Verified {
		args := []string{r.GitDir, r.WorkTree, "status", "--porcelain"}
		cmd := exec.Command("git", args...)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Run()
		if str := out.String(); str != "" {
			r.Clean = false
		} else {
			r.Clean = true
		}
	}
}

func (r *Repo) gitLocalBranch() {
	if r.Verified {
		args := []string{r.GitDir, r.WorkTree, "rev-parse", "--abbrev-ref", "HEAD"}
		cmd := exec.Command("git", args...)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Run()
		if str := out.String(); str != "" {
			r.LocalBranch = trim(out.String())
		}
	}
}

func (r *Repo) gitLocalSHA() {
	if r.Verified {
		args := []string{r.GitDir, r.WorkTree, "rev-parse", "@"}
		cmd := exec.Command("git", args...)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Run()
		if str := out.String(); str != "" {
			r.LocalSHA = trim(out.String())
		}
	}
}

func (r *Repo) gitUpstreamSHA() {
	if r.Verified {
		args := []string{r.GitDir, r.WorkTree, "rev-parse", "@{u}"}
		cmd := exec.Command("git", args...)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Run()
		if str := out.String(); str != "" {
			r.UpstreamSHA = trim(out.String())
		}
	}
}

func (r *Repo) gitMergeBaseSHA() {
	if r.Verified {
		args := []string{r.GitDir, r.WorkTree, "merge-base", "@", "@{u}"}
		cmd := exec.Command("git", args...)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Run()
		if str := out.String(); str != "" {
			r.MergeBaseSHA = trim(out.String())
		}
	}
}

func (r *Repo) gitRevParseUpstream() {
	if r.Verified {
		args := []string{r.GitDir, r.WorkTree, "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}"}
		cmd := exec.Command("git", args...)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Run()
		if str := out.String(); str != "" {
			r.UpstreamBranch = trim(out.String())
		} else {
			fmt.Printf("can't get upstream for %v\n ", r.Name)
		}
	}
}

func (r *Repo) gitDiffsNameOnly() {
	if r.Verified {
		args := []string{r.GitDir, r.WorkTree, "diff", "--name-only", "@{u}"}
		cmd := exec.Command("git", args...)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Run()
		if str := out.String(); str != "" {
			r.DiffFiles = strings.Fields(str)
		} else {
			r.DiffFiles = make([]string, 0)
		}
	}
}

func (r *Repo) getDiffSummary() {
	if r.Verified && len(r.DiffFiles) > 0 {
		r.DiffCount = len(r.DiffFiles)
		// var b bytes.Buffer

		// for _, d := range r.DiffFiles {
		// 	ld := len(strings.Join(r.DiffFiles, ", ")) // length of diff string

		// }

		switch {
		case r.DiffCount == 0:
			r.DiffSummary = "" // r.DiffSummary = "No diffs"
			r.DiffStatus = false
		case r.DiffCount == 1:
			r.DiffSummary = fmt.Sprintf(r.DiffFiles[0])
			r.DiffStatus = true
		case r.DiffCount >= 2:
			var b bytes.Buffer
			t := 0
			for _, d := range r.DiffFiles {
				if b.Len() <= 25 {
					d = fmt.Sprintf("%v, ", d)
					b.WriteString(d)
					t++
				} else {
					break
				}
			}
			s := b.String()
			s = strings.TrimSuffix(s, ", ")
			if t != len(r.DiffFiles) {
				s = fmt.Sprintf("%v...", s)
			}
			r.DiffSummary = s
			r.DiffStatus = true
		}
	}
}

func (r *Repo) gitShortstat() {
	if r.Verified {
		args := []string{r.GitDir, r.WorkTree, "diff", "--shortstat"}
		cmd := exec.Command("git", args...)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Run()
		if str := out.String(); str != "" {
			r.ShortStat = trim(str)
		}
	}
}

func (r *Repo) getShortInts() {
	if r.Verified {
		if r.ShortStat != "" {
			rxi := regexp.MustCompile(`changed, (.*)? insertions`)
			rxs := rxi.FindStringSubmatch(r.ShortStat)
			if len(rxs) == 2 {
				s := rxs[1]
				if i, err := strconv.Atoi(s); err == nil {
					r.ShortStatPlus = i // FLAG: r.PlusCount
				} else {
					fmt.Println(err)
				}
			}

			rxd := regexp.MustCompile(`\(\+\), (.*)? deletions`)
			rxs = rxd.FindStringSubmatch(r.ShortStat)
			if len(rxs) == 2 {
				s := rxs[1]
				if i, err := strconv.Atoi(s); err == nil {
					r.ShortStatMinus = i // FLAG: r.MinusCount
				}
			}
		}
	}
}

func (r *Repo) gitUntracked() {
	if r.Verified {
		args := []string{r.GitDir, r.WorkTree, "ls-files", "--others", "--exclude-standard"}
		cmd := exec.Command("git", args...)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Run()
		if str := out.String(); str != "" {
			ufr := strings.Fields(str) // untracked files raw
			for _, f := range ufr {
				f = lastPathSelection(f)
				r.UntrackedFiles = append(r.UntrackedFiles, f)
			}

		} else {
			r.UntrackedFiles = make([]string, 0)
		}
	}
}

func (r *Repo) getUntrackedSummary() {
	if r.Verified {
		r.UntrackedCount = len(r.UntrackedFiles)
		switch {
		case r.UntrackedCount == 0:
			r.UntrackedSummary = "No untracked files"
			r.UntrackedStatus = false
		case r.UntrackedCount == 1:
			r.UntrackedSummary = fmt.Sprintf(r.UntrackedFiles[0])
			r.UntrackedStatus = true
		case r.UntrackedCount >= 2:
			var b bytes.Buffer
			t := 0
			// FLAG: also limit the size of file names?
			for _, d := range r.UntrackedFiles {
				if b.Len() <= 25 {
					d = fmt.Sprintf("%v, ", d)
					b.WriteString(d)
					t++
				} else {
					break
				}
			}
			s := b.String()
			s = strings.TrimSuffix(s, ", ")
			if t != r.UntrackedCount {
				s = fmt.Sprintf("%v...", s)
			}
			r.UntrackedSummary = s
			r.UntrackedStatus = true
		}
	}
}

func (r *Repo) getUpstreamStatus() {
	if r.Verified {
		switch {
		case r.LocalSHA == r.UpstreamSHA:
			r.Upstream = "Up-To-Date"
		case r.LocalSHA == r.MergeBaseSHA:
			r.Upstream = "Behind"
		case r.UpstreamSHA == r.MergeBaseSHA:
			r.Upstream = "Ahead"
		}
	}
}

func (r *Repo) getPhase() {
	if r.Verified {
		switch {
		case (r.Clean == true && r.UntrackedStatus == false && r.Upstream == "Ahead"):
			r.Phase = "Ahead"
		case (r.Clean == true && r.UntrackedStatus == false && r.Upstream == "Behind"):
			r.Phase = "Behind"
		case (r.Clean == false && r.UntrackedStatus == false && r.Upstream == "Up-To-Date"):
			r.Phase = "Dirty"
		case (r.Clean == false && r.UntrackedStatus == true && r.Upstream == "Up-To-Date"):
			r.Phase = "DirtyUntracked"
		case (r.Clean == false && r.UntrackedStatus == false && r.Upstream == "Ahead"):
			r.Phase = "DirtyAhead"
		case (r.Clean == false && r.UntrackedStatus == false && r.Upstream == "Behind"):
			r.Phase = "DirtyBehind"
		case (r.Clean == false && r.UntrackedStatus == true && r.Upstream == "Up-To-Date"):
			r.Phase = "Untracked"
		case (r.Clean == false && r.UntrackedStatus == true && r.Upstream == "Ahead"):
			r.Phase = "UntrackedAhead"
		case (r.Clean == false && r.UntrackedStatus == true && r.Upstream == "Behind"):
			r.Phase = "UntrackedBehind"
		case (r.Clean == true && r.UntrackedStatus == false && r.Upstream == "Up-To-Date"):
			r.Phase = "Up-To-Date"
		default:
			r.Phase = "wtf"
			fmt.Printf("%v %v %v", r.Clean, r.UntrackedStatus, r.Upstream)
		}
	}
}

// FLAG: overhaul with bytes.buffer at somepoint

func (r *Repo) gitPull(e Emoji, f Flags) {
	targetPrint(f, "%v %v pulling from %v", e.Ship, r.Name, r.Remote)

	args := []string{"-C", r.Path, "pull"}
	cmd := exec.Command("git", args...)
	cmd.Run()
}

func (r *Repo) gitPush(e Emoji, f Flags) {
	targetPrint(f, "%v %v pushing to %v", e.Rocket, r.Name, r.Remote)

	args := []string{"-C", r.Path, "push"}
	cmd := exec.Command("git", args...)
	cmd.Run()
}

func (r *Repo) gitAdd(e Emoji, f Flags) {
	targetPrint(f, "%v  %v adding changes [%v]{%v}(+%v|-%v)", e.Satellite, r.Name, r.DiffCount, r.DiffSummary, r.ShortStatPlus, r.ShortStatMinus)
	args := []string{"-C", r.Path, "add", "-A"}
	cmd := exec.Command("git", args...)
	cmd.Run()
}

// FLAG: check for stash?

func (r *Repo) gitCommit(e Emoji, f Flags) {
	targetPrint(f, "%v %v committing changes [%v]{%v}(+%v|-%v)", e.Fire, r.Name, r.DiffCount, r.DiffSummary, r.ShortStatPlus, r.ShortStatMinus)

	args := []string{"-C", r.Path, "commit", "-m", r.GitMessage}
	cmd := exec.Command("git", args...)
	cmd.Run()
}

func (r *Repo) gitStash(e Emoji, f Flags) {
	targetPrint(f, "%v  %v stashing changes", e.Squirrel, r.Name)

}

func (r *Repo) gitPop(e Emoji, f Flags) {
	targetPrint(f, "%v %v popping changes", e.Popcorn, r.Name)
}

// func (r *Repo) run(e Emoji, f Flags) {
// 	switch r.GitAction {
// 	case "pull":
// 		r.gitPull(e, f)
// 	case "push":
// 		r.gitPush(e, f)
// 	case "add-commit-push":
// 		r.gitAdd(e, f)
// 		r.gitCommit(e, f)
// 		r.gitPush(e, f)
// 	case "stash-pull-pop-commit-push":
// 		r.gitStash(e, f)
// 		r.gitPull(e, f)
// 		r.gitPop(e, f)
// 		r.gitCommit(e, f)
// 		r.gitPush(e, f)
// 	}
// }

// Returns *Div.
// d.Path set from Zone.Path (zp) and bundle.Division (bd)
// Note: zone and bundle aren't defined structs

func initDiv(zp string, bd string) *Div {
	d := new(Div)
	var b bytes.Buffer
	b.WriteString(validatePath(zp))

	if bd != "main" {
		b.WriteString("/")
		b.WriteString(bd)
	}

	d.Name = bd
	d.Path = b.String()

	return d
}

type Div struct {
	// initDiv
	Name string
	Path string

	// initDivsRepos
	Repos Repos

	// verifyDivs
	ReposSummary string
	PathCreated  bool
	PathVerified bool
	Summary      string
	PathError    string

	// setActions
	CompleteRepos  Repos
	PendingRepos   Repos
	ScheduledRepos Repos
	SkippedRepos   Repos
}

func (d *Div) verify(e Emoji, f Flags, w *Workspace) {
	// check path; create if missing and active
	_, err := os.Stat(d.Path)
	if os.IsNotExist(err) && isActive(f) {
		targetPrint(f, "%v creating %v", e.Folder, d.Path)

		os.MkdirAll(d.Path, 0777)
		d.PathCreated = true
		w.CreatedDivs = append(w.CreatedDivs, d)
	}

	// check path
	info, err := os.Stat(d.Path)

	switch {
	case err != nil:
		d.PathVerified = false
		d.PathError = "No directory"
	case os.IsNotExist(err):
		d.PathVerified = false
		d.PathError = "No directory"
	case !info.IsDir():
		d.PathVerified = false
		d.PathError = "File blocking path"
	case noPermission(info):
		d.PathVerified = false
		d.PathError = "No permission"
	default:
		d.PathVerified = true
		d.PathError = ""
	}

	if d.PathVerified {
		w.VerifiedDivs = append(w.VerifiedDivs, d)
		d.Summary = fmt.Sprintf("%v %v", e.Checkmark, d.Path)
	} else {
		d.Summary = fmt.Sprintf("%v %v (%v)", e.Slash, d.Path, d.PathError)
	}
}

// div fns here

func (d *Div) getReposSummary() {
	switch len(d.Repos) {
	case 0:
		d.ReposSummary = "N/A"
	case 1:
		d.ReposSummary = d.Repos[0].Name
	default:
		var sl []string

		for _, r := range d.Repos {
			if lr := len(strings.Join(sl, ", ")); lr <= 22 {
				if len(r.Name) <= 12 {
					sl = append(sl, r.Name)
				}
			}
		}

		s := strings.Join(sl, ", ")

		if len(sl) != len(d.Repos) {
			s = fmt.Sprintf("%v...", s)
		}
		d.ReposSummary = s
	}
}

func (d *Div) getPendingReposSummary() {
	var sl []string

	switch len(d.PendingRepos) {
	case 0:
		d.Summary = "N/A"
	case 1:
		d.Summary = d.PendingRepos[0].Name
	default:
		for _, r := range d.PendingRepos {
			if lsl := len(strings.Join(sl, ", ")); lsl <= 20 {
				// fmt.Println(r.Name)
				sl = append(sl, r.Name)
			}
		}

		sj := strings.Join(sl, ", ")

		if len(sl) != len(d.PendingRepos) {
			var lb bytes.Buffer
			sj = strings.TrimSuffix(sj, ", ")
			lb.WriteString(sj)
			lb.WriteString(",...")
			d.Summary = lb.String()
		} else {
			d.Summary = sj
		}
	}
}

// divs and repos sorting

type Repos []*Repo

func (rs Repos) verify(e Emoji, f Flags, w *Workspace) {
	var wg sync.WaitGroup
	for i := range rs {
		wg.Add(1)
		go func(r *Repo) {
			defer wg.Done()
			r.verify(e, f, w)
		}(rs[i])
	}
	wg.Wait()
}

func (rs Repos) summary(e Emoji, f Flags, t *Timer, w *Workspace) {

	// timer
	// t.ReposParse = mark(t.Start)
	// t.ReposSplit = t.ReposParse - t.DivsReposParse

	// summary
	var b bytes.Buffer

	if len(w.VerifiedRepos) == len(rs) {
		b.WriteString(e.ThumbsUp)
	} else {
		b.WriteString(e.Warning)
	}

	b.WriteString(" [")
	b.WriteString(strconv.Itoa(len(w.VerifiedRepos)))
	b.WriteString("/")
	b.WriteString(strconv.Itoa(len(rs)))
	b.WriteString("] repos up to date")

	if len(w.ClonedRepos) >= 1 {
		b.WriteString(", cloned (")
		b.WriteString(strconv.Itoa(len(w.ClonedRepos)))
		b.WriteString(")")
	}

	b.WriteString(" {")
	// b.WriteString(t.ReposSplit.String())
	b.WriteString(" / ")
	// b.WriteString(t.ReposParse.String())
	b.WriteString("}")

	// print
	targetPrint(f, b.String())
}

func (rs Repos) itemIndex(r *Repo) int {
	for i, rl := range rs {
		if rl.Name == r.Name {
			return i
		}
	}
	return -1
}

type Divs []*Div

func (dvs Divs) verify(e Emoji, f Flags, w *Workspace) {
	for _, d := range dvs {
		d.verify(e, f, w)
	}
}

func (dvs Divs) summary(e Emoji, f Flags, t *Timer, w *Workspace) {

	// timer
	// t.DivsParse = mark(t.Start)
	// t.DivsSplit = t.DivsParse - t.DivsReposParse

	// summary (div)
	for _, d := range dvs {
		targetPrint(f, d.Summary)
	}

	// summary (divs)
	var b bytes.Buffer

	if len(w.VerifiedDivs) == len(dvs) {
		b.WriteString(e.ThumbsUp)
	} else {
		b.WriteString(e.Slash)
	}

	b.WriteString(" [")
	b.WriteString(strconv.Itoa(len(w.VerifiedDivs)))
	b.WriteString("/")
	b.WriteString(strconv.Itoa(len(dvs)))
	b.WriteString("] divs verified")

	if len(w.CreatedDivs) >= 1 {
		b.WriteString(", created (")
		b.WriteString(strconv.Itoa(len(w.CreatedDivs)))
		b.WriteString(")")
	}

	b.WriteString(" {")
	// b.WriteString(t.DivsSplit.String())
	b.WriteString(" / ")
	// b.WriteString(t.DivsParse.String())
	b.WriteString("}")

	// print
	targetPrint(f, b.String())
}

func (dvs Divs) dispatch(e Emoji, f Flags, w *Workspace) {
	for _, d := range dvs {
		d.dispatch(e, f, w)
	}
}

func (d *Div) dispatch(e Emoji, f Flags, w *Workspace) {
	// divs.sort() -------------------------------------
	for _, r := range d.Repos {
		switch {
		case r.Phase == "Ahead" && f.Mode == "logout":
			r.Status = "Scheduled"
			r.GitAction = "push"
			d.ScheduledRepos = append(d.ScheduledRepos, r)
		case r.Phase == "Behind" && f.Mode == "login":
			r.Status = "Scheduled"
			r.GitAction = "pull"
			d.ScheduledRepos = append(d.ScheduledRepos, r)
		case r.Phase == "Up-To-Date" && r.InfoVerified == true:
			r.Status = "Complete"
			d.CompleteRepos = append(d.CompleteRepos, r)
		case r.InfoVerified == false:
			r.Status = "Skipped"
			d.SkippedRepos = append(d.SkippedRepos, r)
		default:
			r.Status = "Pending"
			d.PendingRepos = append(d.PendingRepos, r)
		}
	}
	// ---------------------------------------------

	var b bytes.Buffer

	if len(d.Repos) == len(d.CompleteRepos) {
		b.WriteString(e.Checkmark)
	} else {
		b.WriteString(e.Warning)
	}

	b.WriteString(" ")
	b.WriteString(d.Path)

	if len(d.Repos) == len(d.CompleteRepos) {
		d.getReposSummary()
		b.WriteString(" [")
		b.WriteString(strconv.Itoa(len(d.CompleteRepos)))
		b.WriteString("/")
		b.WriteString(strconv.Itoa(len(d.Repos)))
		b.WriteString("]{")
		b.WriteString(d.ReposSummary)
		b.WriteString("}")

	} else {
		d.getPendingReposSummary()
		b.WriteString(" (")
		b.WriteString(strconv.Itoa(len(d.PendingRepos)))
		b.WriteString("){")
		b.WriteString(d.Summary)
		b.WriteString("}")
	}

	// print summary
	targetPrint(f, b.String())

	// loop over pending, confirm action
	for _, r := range d.PendingRepos {
		var b bytes.Buffer

		switch r.Phase {
		case "Ahead":
			b.WriteString(e.Bunny)
			b.WriteString(" ")
			b.WriteString(r.Name)
			b.WriteString(" is ahead of ")
			b.WriteString(r.UpstreamBranch)
		case "Behind":
			b.WriteString(e.Turtle)
			b.WriteString(" ")
			b.WriteString(r.Name)
			b.WriteString(" is behind")
			b.WriteString(r.UpstreamBranch)
		case "Dirty", "DirtyUntracked", "DirtyAhead", "DirtyBehind":
			b.WriteString(e.Pig)
			b.WriteString(" ")
			b.WriteString(r.Name)
			b.WriteString(" is dirty [")
			b.WriteString(strconv.Itoa((r.DiffCount)))
			b.WriteString("]{")
			b.WriteString(r.DiffSummary)
			b.WriteString("}(+")
			b.WriteString(strconv.Itoa(r.ShortStatPlus))
			b.WriteString("|-")
			b.WriteString(strconv.Itoa(r.ShortStatMinus))
			b.WriteString(")")
		case "Untracked", "UntrackedAhead", "UntrackedBehind":
			b.WriteString(e.Pig)
			b.WriteString(" ")
			b.WriteString(r.Name)
			b.WriteString(" has untracked files[")
			b.WriteString(strconv.Itoa(r.UntrackedCount))
			b.WriteString("]{")
			b.WriteString(r.UntrackedSummary)
			b.WriteString("}")
		case "Up-To-Date":
			b.WriteString(e.Checkmark)
			b.WriteString(" ")
			b.WriteString(r.Name)
			b.WriteString(" is up to date with ")
			b.WriteString(r.UpstreamBranch)
		}

		switch r.Phase {
		case "DirtyUntracked":
			b.WriteString(" with untracked files [")
			b.WriteString(strconv.Itoa(r.UntrackedCount))
			b.WriteString("]{")
			b.WriteString(r.UntrackedSummary)
			b.WriteString("}")
		case "DirtyAhead":
			b.WriteString(" & ahead of ")
			b.WriteString(r.UpstreamBranch)
		case "DirtyBehind":
			b.WriteString(" & behind")
			b.WriteString(r.UpstreamBranch)
		case "UntrackedAhead":
			b.WriteString(" & is ahead of ")
			b.WriteString(r.UpstreamBranch)
		case "UntrackedBehind":
			b.WriteString(" & is behind")
			b.WriteString(r.UpstreamBranch)
		}

		// print prompt (part 1)
		targetPrint(f, b.String())

		// print prompt (part 2)
		switch {
		case r.Phase == "Ahead":
			r.GitAction = "push"
			fmt.Printf("%v push changes to %v? ", e.Rocket, r.Remote)
		case r.Phase == "Behind":
			r.GitAction = "pull"
			fmt.Printf("%v pull changes from %v? ", e.Ship, r.Remote)
		case r.Phase == "Dirty":
			r.GitAction = "add-commit-push"
			fmt.Printf("%v add all files, commit and push to %v? ", e.Clipboard, r.Remote)
		case r.Phase == "DirtyUntracked":
			r.GitAction = "add-commit-push"
			fmt.Printf("%v add all files, commit and push to %v? ", e.Clipboard, r.Remote)
		case r.Phase == "DirtyAhead":
			r.GitAction = "add-commit-push"
			fmt.Printf("%v add all files, commit and push to %v? ", e.Clipboard, r.Remote)
		case r.Phase == "DirtyBehind":
			r.GitAction = "stash-pull-pop-commit-push"
			fmt.Printf("%v stash all files, pull changes, commit and push to %v? ", e.Clipboard, r.Remote)
		case r.Phase == "Untracked":
			r.GitAction = "add-commit-push"
			fmt.Printf("%v add all files, commit and push to %v? ", e.Clipboard, r.Remote)
		case r.Phase == "UntrackedAhead":
			r.GitAction = "add-commit-push"
			fmt.Printf("%v add all files, commit and push to %v? ", e.Clipboard, r.Remote)
		case r.Phase == "UntrackedBehind":
			r.GitAction = "stash-pull-pop-commit-push"
			fmt.Printf("%v stash all files, pull changes, commit and push to %v? ", e.Clipboard, r.Remote)
		}

		rdr := bufio.NewReader(os.Stdin)
		in, err := rdr.ReadString('\n')

		if err != nil {
			r.GitConfirmed = false
		} else {
			in = strings.TrimSuffix(in, "\n")
			switch in {
			case "please", "y", "yes", "ys", "1", "ok", "push", "pull", "sure", "you betcha", "do it":
				r.GitConfirmed = true
			case "n", "no", "nah", "0", "stop", "skip", "abort", "halt", "quit":
				r.GitConfirmed = false
			default:
				r.GitConfirmed = false
			}
		}

		if r.GitConfirmed == true && strings.Contains(r.GitAction, "commit") {
			if hasEmoji(f) {
				fmt.Printf("%v commit message: ", e.Memo)
			} else {
				fmt.Printf("commit message: ")
			}

			rdr := bufio.NewReader(os.Stdin)
			in, _ := rdr.ReadString('\n')
			switch in {
			case "n", "no", "nah", "0", "stop", "skip", "abort", "halt", "quit":
				r.GitConfirmed = false
				r.GitMessage = ""
				d.SkippedRepos = append(d.SkippedRepos, r)
			default:
				r.GitConfirmed = true
				r.GitMessage = in
				d.ScheduledRepos = append(d.ScheduledRepos, r)
			}
		} else if r.GitConfirmed == true {
			d.ScheduledRepos = append(d.ScheduledRepos, r)
		} else if r.GitConfirmed == false {
			d.SkippedRepos = append(d.SkippedRepos, r)
		}
	}
	// clear out pending at this point?
}

// sorting

func (rs Repos) Len() int           { return len(rs) }
func (rs Repos) Swap(i, j int)      { rs[i], rs[j] = rs[j], rs[i] }
func (rs Repos) Less(i, j int) bool { return rs[i].Name < rs[j].Name }

func (dvs Divs) Len() int           { return len(dvs) }
func (dvs Divs) Swap(i, j int)      { dvs[i], dvs[j] = dvs[j], dvs[i] }
func (dvs Divs) Less(i, j int) bool { return dvs[i].Path < dvs[j].Path }

// this is the key part of all this...
func initDivsRepos(c Config, e Emoji, f Flags, t *Timer) (dvs Divs, rs Repos) {

	// print
	targetPrint(f, "%v parsing divs|repos", e.Pager)

	// initialize divs and repos from config
	for _, z := range c.Zones {
		for _, bl := range z.Bundles {
			d := initDiv(z.Path, bl.Division)
			for _, rn := range bl.Repos {
				r := initRepo(d, rn, bl.User, bl.Remote, bl.Division)
				rs = append(rs, r)
				d.Repos = append(d.Repos, r)
			}
			dvs = append(dvs, d)
		}
	}

	// sort
	sort.Sort(Divs(dvs))
	sort.Sort(Repos(rs))

	// timer
	t.markMoment("parse-divs-repos")

	// print
	targetPrint(f, "%v [%v|%v] divs|repos {%v / %v}", e.FaxMachine, len(dvs), len(rs), t.getSplit(), t.getTime())

	return dvs, rs
}

// WORKSPACE

func initWorkspace() *Workspace {
	w := new(Workspace)
	return w
}

type Workspace struct {

	// verifyDivs
	CreatedDivs  Divs
	VerifiedDivs Divs

	// verifyRepos
	ClonedRepos   Repos
	VerifiedRepos Repos

	// FLAG: Streamlining setActions
	// setActions

	CompleteDivs   Divs
	IncompleteDivs Divs

	IsComplete     bool
	CompleteRepos  Repos
	PendingRepos   Repos
	ScheduledRepos Repos
	SkippedRepos   Repos
	PendingCount   int
	DrySummary     string

	// getSummary(*)
	PushRepos      Repos
	PushSummary    string
	PullRepos      Repos
	PullSummary    string
	CommitRepos    Repos
	CommitSummary  string
	SkippedSummary string
	// IsApproved      bool
	ChangesApproved bool
}

func (w *Workspace) getSkippedSummary(e Emoji) {
	if len(w.SkippedRepos) > 0 {
		var b bytes.Buffer
		var fs []string // file slice
		var ds []string // diff slice
		var dc int      // diff count
		var pc int      // plus count
		var mc int      // minus count
		var fj string   // file join
		var dj string   // diff join

		for _, r := range w.SkippedRepos {

			lr := len(strings.Join(fs, ", ")) // length (of) repo string

			if lr <= 20 {
				fs = append(fs, r.Name)
			}

			for _, d := range r.DiffFiles {
				ld := len(strings.Join(ds, ", ")) // length (of) div string

				if ld <= 20 {
					ds = append(ds, d)
				}

			}

			// count of Diffs, Additions (Plus) and Subtractions (Minus)

			dc += r.DiffCount
			pc += r.ShortStatPlus
			mc += r.ShortStatMinus
		}

		fj = strings.Join(fs, ", ")
		dj = strings.Join(ds, ", ")

		// FLAG: logic here?
		if len(fs) != len(w.SkippedRepos) {
			var lb bytes.Buffer
			fj = strings.TrimSuffix(fj, ", ")
			lb.WriteString(fj)
			lb.WriteString(",...")
			fj = lb.String()
		}

		if len(ds) != dc {
			var lb bytes.Buffer
			dj = strings.TrimSuffix(dj, ", ")
			lb.WriteString(dj)
			lb.WriteString(",...")
			dj = lb.String()
		}

		b.WriteString(e.Slash)
		b.WriteString(" skipped = [")
		b.WriteString(strconv.Itoa(len(w.SkippedRepos)))
		b.WriteString("]{")
		b.WriteString(fj)
		b.WriteString("} [")
		b.WriteString(strconv.Itoa(dc))
		b.WriteString("]{")
		b.WriteString(dj)
		b.WriteString("}")
		w.SkippedSummary = b.String()
	}
}

// FLAG: build in support for untracked files too
func (w *Workspace) getCommitSummary(e Emoji) {
	if len(w.CommitRepos) > 0 {
		var b bytes.Buffer
		var fs []string // file slice
		var ds []string // diff slice
		var us []string // untracked slice
		var dc int      // diff count
		var pc int      // plus count
		var mc int      // minus count
		var uc int      // untracked count
		var fj string   // file join
		var dj string   // diff join
		var uj string   // untracked join

		for _, r := range w.CommitRepos {
			lr := len(strings.Join(fs, ", ")) // length (of) repo string
			if lr <= 12 {
				fs = append(fs, r.Name)
			}

			for _, d := range r.DiffFiles {
				ld := len(strings.Join(ds, ", ")) // length (of) div string
				if ld <= 12 {
					ds = append(ds, d)
				}

			}

			for _, f := range r.UntrackedFiles {
				lu := len(strings.Join(us, ", "))
				if lu <= 12 {
					us = append(us, f)
				}
			}

			// count of Diffs, Additions (Plus) and Subtractions (Minus)

			dc += r.DiffCount
			pc += r.ShortStatPlus
			mc += r.ShortStatMinus
			uc += r.UntrackedCount
		}

		fj = strings.Join(fs, ", ")
		dj = strings.Join(ds, ", ")
		uj = strings.Join(us, ", ")

		if len(fs) != len(w.CommitRepos) {
			var lb bytes.Buffer
			fj = strings.TrimSuffix(fj, ", ")
			lb.WriteString(fj)
			lb.WriteString(",...")
			fj = lb.String()
		}

		if len(ds) != dc {
			var lb bytes.Buffer
			dj = strings.TrimSuffix(dj, ", ")
			lb.WriteString(dj)
			lb.WriteString(",...")
			dj = lb.String()
		}

		if len(us) != uc {
			var lb bytes.Buffer
			lb.WriteString(uj)
			lb.WriteString(",...")
			uj = lb.String()
		}

		b.WriteString(e.Clipboard)
		b.WriteString(" commit ")

		if w.SkippedSummary != "" {
			b.WriteString(" ")
		}

		b.WriteString("= [")
		b.WriteString(strconv.Itoa(len(w.CommitRepos)))
		b.WriteString("]{")
		b.WriteString(fj)
		b.WriteString("} [")
		b.WriteString(strconv.Itoa(dc))
		b.WriteString("]{")
		b.WriteString(dj)
		b.WriteString("}(+")
		b.WriteString(strconv.Itoa(pc))
		b.WriteString("/-")
		b.WriteString(strconv.Itoa(mc))
		b.WriteString(")")

		if uc > 0 {
			b.WriteString(" [")
			b.WriteString(strconv.Itoa(uc))
			b.WriteString("]{")
			b.WriteString(uj)
			b.WriteString("}")
		}

		w.CommitSummary = b.String()
	}
}

func (w *Workspace) getPushSummary(e Emoji) {
	if len(w.PushRepos) > 0 {
		var b bytes.Buffer
		var fs []string // file slice
		var ds []string // diff slice
		var dc int      // diff count
		var fj string   // file join
		var dj string   // diff join

		for _, r := range w.PushRepos {

			lr := len(strings.Join(fs, ", ")) // length (of) repo string

			if lr <= 20 {
				fs = append(fs, r.Name)
			}

			for _, d := range r.DiffFiles {
				ld := len(strings.Join(ds, ", ")) // length (of) div string

				if ld <= 20 {
					ds = append(ds, d)
				}

			}

			// count of Diffs, Additions (Plus) and Subtractions (Minus)

			dc += r.DiffCount
		}

		fj = strings.Join(fs, ", ")
		dj = strings.Join(ds, ", ")

		if len(fs) != len(w.PushRepos) {
			var lb bytes.Buffer
			fj = strings.TrimSuffix(fj, ", ")
			lb.WriteString(fj)
			lb.WriteString(",...")
			fj = lb.String()
		}

		if len(ds) != dc {
			var lb bytes.Buffer
			dj = strings.TrimSuffix(dj, ", ")
			lb.WriteString(dj)
			lb.WriteString(",...")
			dj = lb.String()
		}

		b.WriteString(e.Rocket)
		b.WriteString(" push ")

		// fmt.Printf("%v | %v\n", len(w.SkippedSummary), len(w.CommitSummary))

		switch {
		case w.SkippedSummary != "" && w.CommitSummary == "":
			b.WriteString("   ")
		case w.SkippedSummary != "" && w.CommitSummary != "":
			b.WriteString("   ")
		case w.CommitSummary != "" && w.SkippedSummary == "":
			b.WriteString("  ")
		}

		b.WriteString("= [")
		b.WriteString(strconv.Itoa(len(w.PushRepos)))
		b.WriteString("]{")
		b.WriteString(fj)
		b.WriteString("} [")
		b.WriteString(strconv.Itoa(dc))
		b.WriteString("]{")
		b.WriteString(dj)
		b.WriteString("}")

		w.PushSummary = b.String()
	}
}

func (w *Workspace) getPullSummary(e Emoji) {
	if len(w.PullRepos) > 0 {
		var b bytes.Buffer
		var fs []string // file slice
		var ds []string // diff slice
		var dc int      // diff count
		var pc int      // plus count
		var mc int      // minus count
		var fj string   // file join
		var dj string   // diff join

		for _, r := range w.PullRepos {

			lr := len(strings.Join(fs, ", ")) // length (of) repo string

			if lr <= 20 {
				fs = append(fs, r.Name)
			}

			for _, d := range r.DiffFiles {
				ld := len(strings.Join(ds, ", ")) // length (of) div string

				if ld <= 20 {
					ds = append(ds, d)
				}

			}

			// count of Diffs, Additions (Plus) and Subtractions (Minus)

			dc += r.DiffCount
			pc += r.ShortStatPlus
			mc += r.ShortStatMinus
		}

		fj = strings.Join(fs, ", ")
		dj = strings.Join(ds, ", ")

		if len(fs) != len(w.PullRepos) {
			var lb bytes.Buffer
			fj = strings.TrimSuffix(fj, ", ")
			lb.WriteString(fj)
			lb.WriteString(",...")
			fj = lb.String()
		}

		if len(ds) != dc {
			var lb bytes.Buffer
			dj = strings.TrimSuffix(dj, ", ")
			lb.WriteString(dj)
			lb.WriteString(",...")
			dj = lb.String()
		}

		b.WriteString(e.Ship)
		b.WriteString(" pull ")

		switch {
		case w.SkippedSummary != "" && w.CommitSummary == "":
			b.WriteString("   ")
		case w.SkippedSummary != "" && w.CommitSummary != "":
			b.WriteString("   ")
		case w.CommitSummary != "" && w.SkippedSummary == "":
			b.WriteString("  ")
		}

		b.WriteString("= [")
		b.WriteString(strconv.Itoa(len(w.PullRepos)))
		b.WriteString("]{")
		b.WriteString(fj)
		b.WriteString("} [")
		b.WriteString(strconv.Itoa(dc))
		b.WriteString("]{")
		b.WriteString(dj)
		b.WriteString("}")
		w.PullSummary = b.String()
	}
}

// FLAG: Double tap on the marking

// FLAG: more Go like to return err/int?

func (w *Workspace) clearPending(r *Repo) {
	i := w.PendingRepos.itemIndex(r)

	if i >= 0 {
		copy(w.PendingRepos[i:], w.PendingRepos[i+1:])
		w.PendingRepos[len(w.PendingRepos)-1] = nil
		w.PendingRepos = w.PendingRepos[:len(w.PendingRepos)-1]
	}

	i = r.Div.PendingRepos.itemIndex(r)

	if i >= 0 {
		copy(r.Div.PendingRepos[i:], r.Div.PendingRepos[i+1:])
		r.Div.PendingRepos[len(r.Div.PendingRepos)-1] = nil
		r.Div.PendingRepos = r.Div.PendingRepos[:len(r.Div.PendingRepos)-1]
	}
}

func (w *Workspace) clearAllPending() {
	for _, r := range w.ScheduledRepos {
		i := w.PendingRepos.itemIndex(r)
		if i >= 0 {
			copy(w.PendingRepos[i:], w.PendingRepos[i+1:])
			w.PendingRepos[len(w.PendingRepos)-1] = nil
			w.PendingRepos = w.PendingRepos[:len(w.PendingRepos)-1]
		}
	}

	for _, r := range w.SkippedRepos {
		i := w.PendingRepos.itemIndex(r)
		if i >= 0 {
			copy(w.PendingRepos[i:], w.PendingRepos[i+1:])
			w.PendingRepos[len(w.PendingRepos)-1] = nil
			w.PendingRepos = w.PendingRepos[:len(w.PendingRepos)-1]
		}
	}
}

func (w *Workspace) dispatch(dvs Divs, e Emoji, f Flags) {

	// compile scheduled, skipped and complete repos
	for _, d := range dvs {
		w.ScheduledRepos = append(w.ScheduledRepos, d.ScheduledRepos...)
		w.SkippedRepos = append(w.SkippedRepos, d.SkippedRepos...)
		w.CompleteRepos = append(w.CompleteRepos, d.CompleteRepos...)
	}

	// FLAG: just put in both?
	if len(w.ScheduledRepos) >= 1 {
		for _, r := range w.ScheduledRepos {
			switch r.GitAction {
			case "push":
				w.PushRepos = append(w.PushRepos, r)
			case "pull":
				w.PullRepos = append(w.PullRepos, r)
			case "add-commit-push":
				w.CommitRepos = append(w.CommitRepos, r)
			case "stash-pull-pop-commit-push":
				w.CommitRepos = append(w.CommitRepos, r)
			}
		}

		// get summaries
		w.getSkippedSummary(e)
		w.getCommitSummary(e)
		w.getPushSummary(e)
		w.getPullSummary(e)

		// FLAG: getSummary should 'hide' the logic

		fmt.Printf("%v [%v] scheduled: \n", e.Unicorn, len(w.ScheduledRepos)) // FLAG: bring back emojiPrint

		if w.PushSummary != "" {
			fmt.Println(w.PushSummary)
		}

		if w.PullSummary != "" {
			fmt.Println(w.PullSummary)
		}

		if w.CommitSummary != "" {
			fmt.Println(w.CommitSummary)
		}

		if w.SkippedSummary != "" {
			fmt.Println(w.SkippedSummary)
		}

		fmt.Printf("%v submit changes? ", e.Traffic)

		rdr := bufio.NewReader(os.Stdin)
		in, err := rdr.ReadString('\n')

		if err != nil {
			w.ChangesApproved = false
		} else {
			in = strings.TrimSuffix(in, "\n")
			switch in {
			case "please", "y", "yes", "ys", "1", "ok", "push", "pull", "sure", "you betcha", "do it":
				w.ChangesApproved = true
			case "n", "no", "nah", "0", "stop", "skip", "abort", "halt", "quit":
				w.ChangesApproved = false
			default:
				w.ChangesApproved = false
			}
		}

		if w.ChangesApproved {
			targetPrint(f, "%v [%v] submitting changes", e.SatelliteDish, len(w.ScheduledRepos))

			var wg sync.WaitGroup
			for i := range w.ScheduledRepos {
				wg.Add(1)
				go func(r *Repo) {
					defer wg.Done()
					switch r.GitAction {
					case "pull":
						r.gitPull(e, f)
					case "push":
						r.gitPush(e, f)
					case "add-commit-push":
						r.gitAdd(e, f)
						r.gitCommit(e, f)
						r.gitPush(e, f)
					case "stash-pull-pop-commit-push":
						r.gitStash(e, f)
						r.gitPull(e, f)
						r.gitPop(e, f)
						r.gitCommit(e, f)
						r.gitPush(e, f)
					}
				}(w.ScheduledRepos[i])
			}
			wg.Wait()
		}

	}

	// sort scheduled into push, pull or commit
}

// FLAG: different? short summary (multi lines) then show each push/pull/commit (multiline)?
// this needs to also check if Skipped and no pending
// one line per psh/pll/cmt

// Utility functions. Repackage and clarify someday.

func clearScreen(f Flags) {
	if isClear(f) || hasEmoji(f) {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func noPermission(info os.FileInfo) bool {
	s := info.Mode().String()[1:4]
	if s != "rwx" {
		return true
	} else {
		return false
	}
}

func validatePath(p string) string {
	if t := strings.TrimPrefix(p, "~/"); t != p {
		u, err := user.Current()

		if err != nil {
			log.Fatalf("Unable to identify the current user")
		}

		t := strings.Join([]string{u.HomeDir, "/", t}, "")
		return strings.TrimSuffix(t, "/")
	}
	return strings.TrimSuffix(p, "/")
}

func lastPathSelection(p string) string {
	if strings.Contains(p, "/") == true {
		sp := strings.SplitAfter(p, "/") // split path
		lp := sp[len(sp)-1]              // last path
		return lp
	} else {
		return p
	}
}

func trim(s string) string {
	return strings.TrimSuffix(s, "\n")
}

func targetPrint(f Flags, s string, z ...interface{}) {
	var p string
	switch {
	case oneLine(f):
	case isVerbose(f) && hasEmoji(f):
		p = fmt.Sprintf(s, z...)
		fmt.Println(p)
	case isVerbose(f) && noEmoji(f):
		p = fmt.Sprintf(s, z...)
		p = strings.TrimPrefix(p, " ")
		p = strings.TrimPrefix(p, " ")
		fmt.Println(p)
	}
}

// --> Main functions

func initRun() (e Emoji, f Flags, rs Repos, dvs Divs, t *Timer, w *Workspace) {

	// initialize Timer, Flags and Emoji
	t = initTimer()
	f = initFlags(e, t)
	e = initEmoji(f, t)

	// clear screen, early messaging
	initPrint(e, f, t)

	// read ~/.gisrc.json, initialize Config
	c := initConfig(e, f, t, w)

	w = initWorkspace()
	dvs, rs = initDivsRepos(c, e, f, t)

	return e, f, rs, dvs, t, w
}

func verifyDivs(e Emoji, f Flags, rs Repos, dvs Divs, t *Timer, w *Workspace) {

	// print
	targetPrint(f, "%v  verifying divs [%v]", e.FileCabinet, len(dvs))

	// verify
	dvs.verify(e, f, w)

	// summary
	dvs.summary(e, f, t, w)
}

func verifyRepos(e Emoji, f Flags, rs Repos, dvs Divs, t *Timer, w *Workspace) {

	// print
	targetPrint(f, "%v verifying repos [%v]", e.Truck, len(rs))

	// verify
	rs.verify(e, f, w)

	// summary
	rs.summary(e, f, t, w)
}

func verifyChanges(e Emoji, f Flags, rs Repos, dvs Divs, t *Timer, w *Workspace) {
	// summarize all repos by div, prompt user for approval and/or commit message
	dvs.dispatch(e, f, w)

	// summarize, prompt user for global approval, submit changes
	w.dispatch(dvs, e, f)
}

func terminateRun(e Emoji, f Flags, rs Repos, dvs Divs, t *Timer, w *Workspace) {
}

func main() {
	e, f, rs, dvs, t, w := initRun()
	verifyDivs(e, f, rs, dvs, t, w)
	verifyRepos(e, f, rs, dvs, t, w)
	verifyChanges(e, f, rs, dvs, t, w) // setActions()
	terminateRun(e, f, rs, dvs, t, w)
}
