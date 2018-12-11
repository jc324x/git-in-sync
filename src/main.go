package main

import (
	// "bufio"
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
	// "sync"
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
	t.markMoment("init-emoji")

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
	t.markMoment("init-flags")

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

// initPrint prints info for Emoji and Flag values.
func initPrint(e Emoji, f Flags, t *Timer) {

	// clears the screen if f.Clear or f.Emoji are true
	clearScreen(f)

	// targetPrint prints a message with or without an emoji if f.Emoji is true or false.
	targetPrint(f, "%v start", e.Clapper)

	// dry run only messaging
	if isDry(f) {
		targetPrint(f, "%v  dry run; no changes will be made", e.Desert)
	}

	// print flag init
	if ft, err := t.getMoment("flags"); err == nil {
		targetPrint(f, "%v parsing flags", e.FlagInHole)
		targetPrint(f, "%v [%v] flags (%v) {%v / %v}", e.Flag, f.Count, f.Summary, ft.Split, ft.Start)
	}

	// print emoji init
	if et, err := t.getMoment("emoji"); err == nil {
		targetPrint(f, "%v initializing emoji", e.CrystalBall)
		targetPrint(f, "%v [%v] emoji {%v / %v}", e.DirectHit, e.Count, et.Split, et.Start)
	}
}

// --> Config: ~/.gisrc.json unmarshalled

type Config struct {
	Bundles []struct {
		Path  string `json:"path"`
		Zones []struct {
			User     string   `json:"user"`
			Remote   string   `json:"remote"`
			Division string   `json:"division"`
			Repos    []string `json:"repositories"`
		} `json:"zones"`
	} `json:"bundles"`
}

// initConfig returns data from ~/.gisrc.json as a Config struct.
func initConfig(e Emoji, f Flags, t *Timer) (c Config) {

	// get the current user, otherwise fatal
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
	err = json.Unmarshal(r, &c)

	if err != nil {
		log.Fatalf("Can't unmarshal JSON from %v\n", g)
	}

	// timer
	t.markMoment("init-config")

	// print
	targetPrint(f, "%v read %v {%v / %v}", e.Book, g, t.getSplit(), t.getTime())

	return c
}

// --> Repo: Repository configuration and information

type Repo struct {

	// initRun -> initRepos -> initRepo
	BundlePath   string // "~/dev"
	ZoneDivision string // "main" or "go-lang"
	ZoneUser     string // "jychri"
	ZoneRemote   string // "github" or "gitlab"
	RepoName     string // "git-in-sync"
	DivPath      string // "/Users/jychri/dev/go-lang/"
	RepoPath     string // "/Users/jychri/dev/go-lang/git-in-sync"
	GitPath      string // "/Users/jychri/dev/go-lang/git-in-sync/.git"
	GitDir       string // "--git-dir=/Users/jychri/dev/go-lang/git-in-sync/.git"
	WorkTree     string // "--work-tree=/Users/jychri/dev/go-lang/git-in-sync"
	RepoURL      string // "https://github.com/jychri/git-in-sync"

	// rs.verifyDivs
	DivPathVerified bool   // true if DivPath verified
	DivPathError    string // error if DivPathVerified is false

	// rs.verifyRepos

	// --- maybe maybe not? we'll see... ---

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

// initRepo returns a *Repo with initial values set.

func initRepo(zd string, zu string, zr string, bp string, rn string) *Repo {

	r := new(Repo)

	// "~/dev", (b)undle(p)ath
	r.BundlePath = bp

	// "main" or "go-lang", (z)one(d)ivision
	r.ZoneDivision = zd

	// "jychri", (z)one(u)ser
	r.ZoneUser = zu

	// "github" or "gitlab", (z)one(r)emote
	r.ZoneRemote = zr

	// "git-in-sync", (r)epo(n)ame
	r.RepoName = rn

	var b bytes.Buffer

	// "/Users/jychri/dev/go-lang/"
	b.WriteString(validatePath(r.BundlePath))
	if r.ZoneDivision != "main" {
		b.WriteString("/")
		b.WriteString(r.ZoneDivision)
	}
	r.DivPath = b.String()

	// "/Users/jychri/dev/go-lang/git-in-sync/"
	b.Reset()
	b.WriteString(r.DivPath)
	b.WriteString("/")
	b.WriteString(r.RepoName)
	r.RepoPath = b.String()

	// "/Users/jychri/dev/go-lang/git-in-sync/.git"
	b.Reset()
	b.WriteString(r.RepoPath)
	b.WriteString("/.git")
	r.GitPath = b.String()

	// "--git-dir=/Users/jychri/dev/go-lang/git-in-sync/.git"
	b.Reset()
	b.WriteString("--git-dir=")
	b.WriteString(r.GitPath)
	r.GitDir = b.String()

	// "--work-tree=/Users/jychri/dev/go-lang/git-in-sync"
	b.Reset()
	b.WriteString("--work-tree=")
	b.WriteString(r.RepoPath)
	r.WorkTree = b.String()

	// "https://github.com/jychri/git-in-sync"
	b.Reset()
	switch r.ZoneRemote {
	case "github":
		b.WriteString("https://github.com/")
	case "gitlab":
		b.WriteString("https://gitlab.com/")
	}
	b.WriteString(r.ZoneUser)
	b.WriteString("/")
	b.WriteString(r.RepoName)
	r.RepoURL = b.String()

	return r
}

// repo fns here

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

// --> Repos: Collection of Repos

type Repos []*Repo

// sort A-Z by r.RepoName
func (rs Repos) sortByName() {
	sort.SliceStable(rs, func(i, j int) bool { return rs[i].RepoName < rs[j].RepoName })
}

// sort A-Z by r.DivPath, then r.RepoName
func (rs Repos) sortByPath() {
	sort.SliceStable(rs, func(i, j int) bool { return rs[i].RepoName < rs[j].RepoName })
	sort.SliceStable(rs, func(i, j int) bool { return rs[i].DivPath < rs[j].DivPath })
}

func (rs Repos) verifyDivs(e Emoji, f Flags, t *Timer) {

	// sort
	rs.sortByPath()

	// total divs
	var dvs []string // divs

	for _, r := range rs {
		dvs = append(dvs, r.DivPath)
	}

	dvs = removeDuplicates(dvs)

	// print
	targetPrint(f, "%v  verifying divs [%v]", e.FileCabinet, len(dvs))

	// created, verified and missing divs
	var cd []string  // created divs
	var vfd []string // verified divs
	var md []string  // missing divs

	for _, r := range rs {

		_, err := os.Stat(r.DivPath)

		if os.IsNotExist(err) && isActive(f) {
			targetPrint(f, "%v creating %v", e.Folder, r.DivPath)
			os.MkdirAll(r.DivPath, 0777)
			cd = append(cd, r.DivPath)
		}

		info, err := os.Stat(r.DivPath)

		switch {
		case noPermission(info):
			r.DivPathVerified = false
			r.DivPathError = "No permission"
			md = append(md, r.DivPath)
		case !info.IsDir():
			r.DivPathVerified = false
			r.DivPathError = "File blocking path"
			md = append(md, r.DivPath)
		case os.IsNotExist(err):
			r.DivPathVerified = false
			r.DivPathError = "No directory"
			md = append(md, r.DivPath)
		case err != nil:
			r.DivPathVerified = false
			r.DivPathError = "No directory"
			md = append(md, r.DivPath)
		default:
			r.DivPathVerified = true
			r.DivPathError = ""
			vfd = append(vfd, r.DivPath)
		}
	}

	// timer
	t.markMoment("verify-divs")

	// remove duplicates from slices
	vfd = removeDuplicates(vfd)
	md = removeDuplicates(md)

	// summary
	var b bytes.Buffer

	if len(dvs) == len(vfd) {
		b.WriteString(e.ThumbsUp)
	} else {
		b.WriteString(e.Slash)
	}

	b.WriteString(" [")
	b.WriteString(strconv.Itoa(len(vfd)))
	b.WriteString("/")
	b.WriteString(strconv.Itoa(len(dvs)))
	b.WriteString("] divs verified")

	if len(cd) >= 1 {
		b.WriteString(", created (")
		b.WriteString(strconv.Itoa(len(cd)))
		b.WriteString(")")
	}

	b.WriteString(" {")
	b.WriteString(t.getSplit().String())
	b.WriteString(" / ")
	b.WriteString(t.getTime().String())
	b.WriteString("}")

	targetPrint(f, b.String())
}

func (rs Repos) verifyRepos(e Emoji, f Flags, t *Timer) {

	// print
	targetPrint(f, "%v verifying repos [%v]", e.Truck, len(rs))
}

func initRepos(c Config, e Emoji, f Flags, t *Timer) (rs Repos) {

	// print
	targetPrint(f, "%v parsing repos", e.Pager)

	// initialize Repos from Config
	for _, bl := range c.Bundles {
		for _, z := range bl.Zones {
			for _, rn := range z.Repos {
				r := initRepo(z.Division, z.User, z.Remote, bl.Path, rn)
				rs = append(rs, r)
			}
		}
	}

	// timer
	t.markMoment("init-repos")

	// print
	targetPrint(f, "%v [%v] repos {%v / %v}", e.FaxMachine, len(rs), t.getSplit(), t.getTime())

	return rs
}

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

func removeDuplicates(ssl []string) (sl []string) {

	smap := make(map[string]bool)

	for i := range ssl {
		if smap[ssl[i]] == true {
		} else {
			smap[ssl[i]] = true
			sl = append(sl, ssl[i])
		}
	}

	return sl
}

// --> Main functions

func initRun() (e Emoji, f Flags, rs Repos, t *Timer) {

	// initialize Timer, Flags and Emoji
	t = initTimer()
	f = initFlags(e, t)
	e = initEmoji(f, t)

	// clear screen, early messaging
	initPrint(e, f, t)

	// read ~/.gisrc.json, initialize Config
	c := initConfig(e, f, t)

	// initialize Repos
	rs = initRepos(c, e, f, t)

	return e, f, rs, t
}

func main() {
	e, f, rs, t := initRun()
	rs.verifyDivs(e, f, t)
	rs.verifyRepos(e, f, t)
	// rs.verifyChanges(e, f, t)
}
