package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html"
	"io"
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
	Briefcase            string
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
	Sheep                string
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
	e.Briefcase = printEmoji(128188)
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
	e.Sheep = printEmoji(128017)
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

// FLAG: silent error / warning flag?
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

// logoutMode returns true if mode is set to "logout"
func logoutMode(f Flags) bool {
	if f.Mode == "logout" {
		return true
	} else {
		return false
	}
}

// loginMode returns true if mode is set to "login"
func loginMode(f Flags) bool {
	if f.Mode == "login" {
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
	if ft, err := t.getMoment("init-flags"); err == nil {
		targetPrint(f, "%v parsing flags", e.FlagInHole)
		targetPrint(f, "%v [%v] flags (%v) {%v / %v}", e.Flag, f.Count, f.Summary, ft.Split, ft.Start)
	}

	// print emoji init
	if et, err := t.getMoment("init-emoji"); err == nil {
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
	BundlePath string // "~/dev"
	Division   string // "main" or "go-lang"
	User       string // "jychri"
	Remote     string // "github" or "gitlab"
	Name       string // "git-in-sync"
	DivPath    string // "/Users/jychri/dev/go-lang/"
	RepoPath   string // "/Users/jychri/dev/go-lang/git-in-sync"
	GitPath    string // "/Users/jychri/dev/go-lang/git-in-sync/.git"
	GitDir     string // "--git-dir=/Users/jychri/dev/go-lang/git-in-sync/.git"
	WorkTree   string // "--work-tree=/Users/jychri/dev/go-lang/git-in-sync"
	URL        string // "https://github.com/jychri/git-in-sync"

	// rs.verifyRepos
	PendingClone bool // true if RepoPath or GitPath are empty

	// rs.verifyDivs, rs.verifyRepos
	Verified     bool   // true if Repo continues to pass verification
	ErrorMessage string // the last error message
	ErrorName    string // name of the last error
	ErrorFirst   string // first line of the last error message
	ErrorShort   string // message in matched short form

	// rs.verifyRepos -> gitVerify -> gitClone
	Cloned bool // true if Repo was cloned

	// rs.verifyRepos -> gitConfigOriginURL
	OriginURL string // "https://github.com/jychri/git-in-sync"

	// rs.verifyRepos -> gitAbbrevRef
	LocalBranch string // `git rev-parse --abbrev-ref HEAD`, "master"

	// rs.verifyRepos -> gitLocalSHA
	LocalSHA string // `git rev-parse @`, "l00000ngSHA1slong324"

	// rs.verifyRepos -> gitUpstreamSHA
	UpstreamSHA string // `git rev-parse @{u}`, "l00000ngSHA1slong324"

	// rs.verifyRepos -> gitMergeBaseSHA
	MergeSHA string // `git merge-base @ @{u}`, "l00000ngSHA1slong324"

	// rs.verifyRepos -> gitRevParseUpstream
	UpstreamBranch string // `git rev-parse --abbrev-ref --symbolic-full-name @{u}`, "..."

	// rs.verifyRepos -> gitDiffsNameOnly
	DiffsNameOnly []string // `git diff --name-only @{u}`, [a, b, c, d, e]
	DiffsSummary  string   // "a, b, c..."

	// rs.verifyRepos -> gitShortstat
	ShortStat        string // `git diff --shortstat`, "x files changed, y insertions(+), z deletions(-)"
	Changed          int    // x
	Insertions       int    // y
	Deletions        int    // z
	ShortStatSummary string // "+y|-z" or "D" for Deleted if (x >= 1 && y == 0 && z == 0)
	Clean            bool   // true if Changed, Insertions and Deletions are all 0

	// rs.verifyRepos -> gitUntracked
	UntrackedFiles   []string // `git ls-files --others --exclude-standard`, [a, b, c, d, e]
	UntrackedSummary string   // "a, b, c..."
	Untracked        bool     // true if if len(r.UntrackedFiles) >= 1

	// rs.verifyRepos -> setStatus
	Category string // Complete, Pending, Skipped, Scheduled
	Status   string // better term?

	// setActions
	GitAction  string // "..."
	GitMessage string //
}

// initRepo returns a *Repo with initial values set.

func initRepo(zd string, zu string, zr string, bp string, rn string) *Repo {

	r := new(Repo)

	// "~/dev", (b)undle(p)ath
	r.BundlePath = bp

	// "main" or "go-lang", (z)one(d)ivision
	r.Division = zd

	// "jychri", (z)one(u)ser
	r.User = zu

	// "github" or "gitlab", (z)one(r)emote
	r.Remote = zr

	// "git-in-sync", (r)epo(n)ame
	r.Name = rn

	var b bytes.Buffer

	// "/Users/jychri/dev/go-lang/"
	b.WriteString(validatePath(r.BundlePath))
	if r.Division != "main" {
		b.WriteString("/")
		b.WriteString(r.Division)
	}
	r.DivPath = b.String()

	// "/Users/jychri/dev/go-lang/git-in-sync/"
	b.Reset()
	b.WriteString(r.DivPath)
	b.WriteString("/")
	b.WriteString(r.Name)
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
	switch r.Remote {
	case "github":
		b.WriteString("https://github.com/")
	case "gitlab":
		b.WriteString("https://gitlab.com/")
	}
	b.WriteString(r.User)
	b.WriteString("/")
	b.WriteString(r.Name)
	r.URL = b.String()

	return r
}

func notVerified(r *Repo) bool {
	if r.Verified == false {
		return true
	} else {
		return false
	}
}

func (r *Repo) markError(e Emoji, f Flags, err string, name string) {
	r.ErrorMessage = err
	r.ErrorName = name
	r.ErrorFirst = firstLine(err)

	if strings.Contains(r.ErrorFirst, "warning") {
		r.Verified = true
	}

	if strings.Contains(r.ErrorFirst, "fatal") {
		r.Verified = false
	}
}

func captureOut(b bytes.Buffer) string {
	return strings.TrimSuffix(b.String(), "\n")
}

func (r *Repo) gitCheckPending(e Emoji, f Flags) {

	// return if not verified
	if notVerified(r) {
		return
	}

	// check if RepoPath and GitPath are accessible
	rinfo, rerr := os.Stat(r.RepoPath)
	ginfo, gerr := os.Stat(r.GitPath)

	switch {
	case isFile(rinfo):
		r.markError(e, f, "fatal: file occupying path", "git-verify")
	case isDirectory(rinfo) && notEmpty(r.RepoPath) && os.IsNotExist(gerr):
		r.markError(e, f, "fatal: directory occupying path", "git-verify")
	case isDirectory(rinfo) && isEmpty(r.RepoPath) && isActive(f):
		r.PendingClone = true
	case os.IsNotExist(rerr) && os.IsNotExist(gerr) && isActive(f):
		r.PendingClone = true
	case isDirectory(rinfo) && isEmpty(r.RepoPath) && isDry(f):
		r.markError(e, f, "fatal: git clone (dry run)", "git-verify")
	case os.IsNotExist(rerr) && os.IsNotExist(gerr) && isActive(f):
		r.markError(e, f, "fatal: git clone (dry run)", "git-verify")
	case isDirectory(rinfo) && isDirectory(ginfo):
		r.Verified = true
	}
}

func (r *Repo) gitClone(e Emoji, f Flags) {

	if r.PendingClone == true {
		// print
		targetPrint(f, "%v cloning %v {%v}", e.Box, r.Name, r.Division)

		// command
		args := []string{"clone", r.URL, r.RepoPath}
		cmd := exec.Command("git", args...)
		var out bytes.Buffer
		var err bytes.Buffer
		cmd.Stderr = &err
		cmd.Stdout = &out
		cmd.Run()

		// check error, set value(s)
		if err := err.String(); err != "" {
			r.markError(e, f, err, "gitClone")
		}

		r.Cloned = true

	}

}

func (r *Repo) gitConfigOriginURL(e Emoji, f Flags) {

	// return if not verified
	if notVerified(r) {
		return
	}

	// command
	args := []string{r.GitDir, "config", "--get", "remote.origin.url"}
	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()

	// trim "\n" from command output
	s := out.String()
	s = strings.TrimSuffix(s, "\n")

	// set OriginURL
	r.OriginURL = s

	// check error, set value(s)
	switch {
	case r.OriginURL == "":
		r.markError(e, f, "fatal: 'origin' does not appear to be a git repository", "gitConfigOriginURL")
	case r.OriginURL != r.URL:
		r.markError(e, f, "fatal: URL != OriginURL", "gitConfigOriginURL")
	}
}

func (r *Repo) gitRemoteUpdate(e Emoji, f Flags) {

	// return if not verified
	if notVerified(r) {
		return
	}

	// command
	args := []string{r.GitDir, r.WorkTree, "fetch", "origin"}
	cmd := exec.Command("git", args...)
	var err bytes.Buffer
	cmd.Stderr = &err
	cmd.Run()

	// Warnings for redirects to "*./git" can be ignored.
	eval := err.String()
	wgit := strings.Join([]string{r.URL}, "/.git") // (w)ith .(git)

	switch {
	case strings.Contains(eval, "warning: redirecting") && strings.Contains(eval, wgit):
		// fmt.Printf("%v - redirect to .git\n", r.Name)
	case eval != "":
		r.markError(e, f, eval, "gitRemoteUpdate")
	}
}

func (r *Repo) gitAbbrevRef(e Emoji, f Flags) {

	// return if not verified
	if notVerified(r) {
		return
	}

	// command
	args := []string{r.GitDir, r.WorkTree, "rev-parse", "--abbrev-ref", "HEAD"}
	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	var err bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err
	cmd.Run()

	// check error, set value(s)
	if err := err.String(); err != "" {
		r.markError(e, f, err, "gitAbbrevRef")
	} else {
		r.LocalBranch = captureOut(out)
	}
}

func (r *Repo) gitLocalSHA(e Emoji, f Flags) {

	// return if not verified
	if notVerified(r) {
		return
	}

	// command
	args := []string{r.GitDir, r.WorkTree, "rev-parse", "@"}
	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	var err bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err
	cmd.Run()

	// check error, set value(s)
	if err := err.String(); err != "" {
		r.markError(e, f, err, "gitLocalSHA")
	} else {
		r.LocalSHA = captureOut(out)
	}
}

func (r *Repo) gitUpstreamSHA(e Emoji, f Flags) {

	// return if not verified
	if notVerified(r) {
		return
	}

	// command
	args := []string{r.GitDir, r.WorkTree, "rev-parse", "@{u}"}
	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	var err bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err
	cmd.Run()

	// check error, set value(s)
	if err := err.String(); err != "" {
		r.markError(e, f, err, "gitUpstreamSHA")
	} else {
		r.UpstreamSHA = captureOut(out)
	}
}

func (r *Repo) gitMergeBaseSHA(e Emoji, f Flags) {

	// return if not verified
	if notVerified(r) {
		return
	}

	// command
	args := []string{r.GitDir, r.WorkTree, "merge-base", "@", "@{u}"}
	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	var err bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err
	cmd.Run()

	// check error, set value(s)
	if err := err.String(); err != "" {
		r.markError(e, f, err, "gitUpstreamSHA")
	} else {
		r.MergeSHA = captureOut(out)
	}
}

func (r *Repo) gitRevParseUpstream(e Emoji, f Flags) {

	// return if not verified
	if notVerified(r) {
		return
	}

	// command
	args := []string{r.GitDir, r.WorkTree, "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}"}
	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	var err bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err
	cmd.Run()

	// check error, set value(s)
	if err := err.String(); err != "" {
		r.markError(e, f, err, "gitRevParseUpstream")
	} else {
		r.UpstreamBranch = captureOut(out)
	}
}

func (r *Repo) gitDiffsNameOnly(e Emoji, f Flags) {

	// return if not verified
	if notVerified(r) {
		return
	}

	// command
	args := []string{r.GitDir, r.WorkTree, "diff", "--name-only", "@{u}"}
	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	var err bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err
	cmd.Run()

	// check error, set value(s)
	if err := err.String(); err != "" {
		r.markError(e, f, err, "gitDiffsNameOnly")
	}

	if str := out.String(); str != "" {
		r.DiffsNameOnly = strings.Fields(str)
		r.DiffsSummary = sliceSummary(r.DiffsNameOnly, 12)
	} else {
		r.DiffsNameOnly = make([]string, 0)
		r.DiffsSummary = ""
	}
}

func (r *Repo) gitShortstat(e Emoji, f Flags) {

	// return if not verified
	if notVerified(r) {
		return
	}

	// command
	args := []string{r.GitDir, r.WorkTree, "diff", "--shortstat"}
	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	var err bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err
	cmd.Run()

	// check error, set value(s)
	if err := err.String(); err != "" {
		r.markError(e, f, err, "gitShortstat")
	} else {
		r.ShortStat = captureOut(out)
	}

	// scrape with regular expressions
	rxc := regexp.MustCompile(`(.*)? file`)
	rxs := rxc.FindStringSubmatch(r.ShortStat)
	if len(rxs) == 2 {
		s := strings.TrimPrefix(rxs[1], " ")
		if i, err := strconv.Atoi(s); err == nil {
			r.Changed = i
		}
	}

	rxi := regexp.MustCompile(`changed, (.*)? insertions`)
	rxs = rxi.FindStringSubmatch(r.ShortStat)
	if len(rxs) == 2 {
		s := rxs[1]
		if i, err := strconv.Atoi(s); err == nil {
			r.Insertions = i
		}
	}

	rxd := regexp.MustCompile(`\(\+\), (.*)? deletions`)
	rxs = rxd.FindStringSubmatch(r.ShortStat)
	if len(rxs) == 2 {
		s := rxs[1]
		if i, err := strconv.Atoi(s); err == nil {
			r.Deletions = i
		}
	}

	// set Rlean and ShortStatSummary
	switch {
	case r.Changed == 0 && r.Insertions == 0 && r.Deletions == 0:
		r.Clean = true
		r.ShortStatSummary = ""
	case r.Changed >= 1 && r.Insertions == 0 && r.Deletions == 0:
		r.Clean = false
		r.ShortStatSummary = ("D")
	default:
		r.Clean = false

		var b bytes.Buffer
		b.WriteString("+")
		b.WriteString(strconv.Itoa(r.Insertions))
		b.WriteString("|-")
		b.WriteString(strconv.Itoa(r.Deletions))
		r.ShortStatSummary = b.String()
	}

	if r.Changed == 0 && r.Insertions == 0 && r.Deletions == 0 {
		r.Clean = true
	} else {
		r.Clean = false
	}

}

func (r *Repo) gitUntracked(e Emoji, f Flags) {

	// return if not verified
	if notVerified(r) {
		return
	}

	// command
	args := []string{r.GitDir, r.WorkTree, "ls-files", "--others", "--exclude-standard"}
	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	var err bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err
	cmd.Run()

	// check error, set value(s)
	if err := err.String(); err != "" {
		r.markError(e, f, err, "gitUntracked")
	}

	if str := out.String(); str != "" {
		ufr := strings.Fields(str) // untracked files raw
		for _, f := range ufr {
			f = lastPathSelection(f)
			r.UntrackedFiles = append(r.UntrackedFiles, f)
			r.UntrackedSummary = sliceSummary(r.UntrackedFiles, 12)
		}
	} else {
		r.UntrackedFiles = make([]string, 0)
	}

	if len(r.UntrackedFiles) >= 1 {
		r.Untracked = true
	}

}

func (r *Repo) setStatus(e Emoji, f Flags) {

	switch {
	case r.LocalSHA == r.UpstreamSHA:
		r.Status = "Up-To-Date"
	case r.LocalSHA == r.MergeSHA:
		r.Status = "Behind"
	case r.UpstreamSHA == r.MergeSHA:
		r.Status = "Ahead"
	}

	switch {
	case r.Verified == false:
		r.Category = "Skipped"
		r.Status = "Error"
	case (r.Clean == true && r.Untracked == false && r.Status == "Ahead"):
		r.Category = "Pending"
		r.Status = "Ahead"
		r.GitAction = "push"
	case (r.Clean == true && r.Untracked == false && r.Status == "Behind"):
		r.Category = "Pending"
		r.Status = "Behind"
		r.GitAction = "pull"
	case (r.Clean == false && r.Untracked == false && r.Status == "Up-To-Date"):
		r.Category = "Pending"
		r.Status = "Dirty"
		r.GitAction = "add-commit-push"
	case (r.Clean == false && r.Untracked == true && r.Status == "Up-To-Date"):
		r.Category = "Pending"
		r.Status = "DirtyUntracked"
		r.GitAction = "add-commit-push"
	case (r.Clean == false && r.Untracked == false && r.Status == "Ahead"):
		r.Category = "Pending"
		r.Status = "DirtyAhead"
		r.GitAction = "add-commit-push"
	case (r.Clean == false && r.Untracked == false && r.Status == "Behind"):
		r.Category = "Pending"
		r.Status = "DirtyBehind"
		r.GitAction = "stash-pull-pop-commit-push"
	case (r.Clean == true && r.Untracked == true && r.Status == "Up-To-Date"):
		r.Category = "Pending"
		r.Status = "Untracked"
		r.GitAction = "add-commit-push"
	case (r.Clean == false && r.Untracked == true && r.Status == "Ahead"):
		r.Category = "Pending"
		r.Status = "UntrackedAhead"
		r.GitAction = "add-commit-push"
	case (r.Clean == false && r.Untracked == true && r.Status == "Behind"):
		r.Category = "Pending"
		r.Status = "UntrackedBehind"
	case (r.Clean == true && r.Untracked == false && r.Status == "Up-To-Date"):
		r.Category = "Complete"
		r.Status = "Up-To-Date"
		r.GitAction = "stash-pull-pop-commit-push"
	default:
		r.Category = "Skipped"
		r.Status = "Unknown"
		r.markError(e, f, "fatal: no matches found in setStatus switch", "setStatus")
	}

	if r.ErrorMessage != "" {
		err := r.ErrorMessage
		switch {
		case strings.Contains(err, "fatal: ambiguous argument 'HEAD'"):
			r.ErrorShort = "fatal: empty repository"
		case strings.Contains(err, "fatal: 'origin' does not appear to be a git repository"):
			r.ErrorShort = "fatal: 'origin' not set"
		case strings.Contains(err, "fatal: URL != OriginURL"):
			r.ErrorShort = "fatal: URL mismatch"
		case strings.Contains(err, "fatal: no matches found"):
			r.ErrorShort = "fatal: no matches found"
		}
	}

	// auto move to scheduled for matching login/logout
	switch {
	case loginMode(f) && r.Category == "Pending" && r.Status == "Behind":
		r.Category = "Scheduled"
	case logoutMode(f) && r.Category == "Pending" && r.Status == "Ahead":
		r.Category = "Scheduled"
	}
}

func (r *Repo) checkConfirmed() {

	// setup reader
	rdr := bufio.NewReader(os.Stdin)
	in, err := rdr.ReadString('\n')

	// return if error
	if err != nil {
		r.Category = "Skipped"
		return
	}

	// trim trailing new line
	in = strings.TrimSuffix(in, "\n")

	switch in {
	case "please", "y", "ye", "yes", "ys", "1", "ok", "push", "pull":
		r.Category = "Scheduled"
	case "you may fire when ready", "do it", "just do it", "you betcha", "sure":
		r.Category = "Scheduled"
	case "n", "no", "nah", "0", "stop", "skip", "abort", "halt", "quit":
		r.Category = "Skipped"
	default:
		r.Category = "Skipped"
	}
}

func (r *Repo) checkCommitMessage() {

	// setup reader
	rdr := bufio.NewReader(os.Stdin)
	in, err := rdr.ReadString('\n')

	// return if error
	if err != nil {
		r.Category = "Skipped"
		return
	}

	// trim trailing new line
	in = strings.TrimSuffix(in, "\n")

	switch in {
	case "n", "no", "nah", "0", "stop", "skip", "abort", "halt", "quit", "exit", "":
		r.Category = "Skipped"
		r.GitMessage = ""
	default:
		r.Category = "Scheduled"
		r.GitMessage = in
	}
}

func (r *Repo) gitPull(e Emoji, f Flags) {
	targetPrint(f, "%v %v pulling from %v", e.Ship, r.Name, r.Remote)

	args := []string{"-C", r.RepoPath, "pull"}
	cmd := exec.Command("git", args...)
	cmd.Run()
}

func (r *Repo) gitPush(e Emoji, f Flags) {
	targetPrint(f, "%v %v pushing to %v", e.Rocket, r.Name, r.Remote)

	args := []string{"-C", r.RepoPath, "push"}
	cmd := exec.Command("git", args...)
	cmd.Run()
}

func (r *Repo) gitAdd(e Emoji, f Flags) {
	targetPrint(f, "%v  %v adding changes [%v]{%v}(%v)", e.Satellite, r.Name, len(r.DiffsNameOnly), r.DiffsSummary, r.ShortStatSummary)
	args := []string{"-C", r.RepoPath, "add", "-A"}
	cmd := exec.Command("git", args...)
	cmd.Run()
}

// FLAG: check for stash?

func (r *Repo) gitCommit(e Emoji, f Flags) {
	targetPrint(f, "%v %v committing changes [%v]{%v}(%v)", e.Fire, r.Name, len(r.DiffsNameOnly), r.DiffsSummary, r.ShortStatSummary)

	args := []string{"-C", r.RepoPath, "commit", "-m", r.GitMessage}
	cmd := exec.Command("git", args...)
	cmd.Run()
}

func (r *Repo) gitStash(e Emoji, f Flags) {
	targetPrint(f, "%v  %v stashing changes", e.Squirrel, r.Name)

}

func (r *Repo) gitPop(e Emoji, f Flags) {
	targetPrint(f, "%v %v popping changes", e.Popcorn, r.Name)
}

// --> Repos: Collection of Repos

type Repos []*Repo

func initRepos(c Config, e Emoji, f Flags, t *Timer) (rs Repos) {

	// print
	targetPrint(f, "%v parsing divs|repos", e.Pager)

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

	// sort
	rs.sortByPath()

	// get all divs, remove duplicates
	var dvs []string // divs

	for _, r := range rs {
		dvs = append(dvs, r.DivPath)
	}

	dvs = removeDuplicates(dvs)

	// print
	targetPrint(f, "%v [%v|%v] divs|repos {%v / %v}", e.FaxMachine, len(dvs), len(rs), t.getSplit(), t.getTime())

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

// sort A-Z by r.Name
func (rs Repos) sortByName() {
	sort.SliceStable(rs, func(i, j int) bool { return rs[i].Name < rs[j].Name })
}

// sort A-Z by r.DivPath, then r.Name
func (rs Repos) sortByPath() {
	sort.SliceStable(rs, func(i, j int) bool { return rs[i].Name < rs[j].Name })
	sort.SliceStable(rs, func(i, j int) bool { return rs[i].DivPath < rs[j].DivPath })
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

	if info == nil {
		return false
	}

	if len(info.Mode().String()) <= 4 {
		return true
	}

	s := info.Mode().String()[1:4]

	if s != "rwx" {
		return true
	} else {
		return false
	}
}

func isDirectory(info os.FileInfo) bool {
	if info == nil {
		return false
	}

	if info.IsDir() {
		return true
	} else {
		return false
	}
}

func isEmpty(p string) bool {
	f, err := os.Open(p)

	if err != nil {
		return false
	}

	_, err = f.Readdir(1)

	if err == io.EOF {
		return true
	}

	return false
}

func notEmpty(p string) bool {
	f, err := os.Open(p)

	if err != nil {
		return false
	}

	_, err = f.Readdir(1)

	if err == io.EOF {
		return false
	}

	return true
}

func isFile(info os.FileInfo) bool {
	if info == nil {
		return false
	}

	if info.IsDir() {
		return false
	} else {
		return true
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

func targetPrint(f Flags, s string, z ...interface{}) {
	var p string // print
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

func firstLine(s string) string {
	lines := strings.Split(strings.TrimSuffix(s, "\n"), "\n")

	if len(lines) >= 1 {
		return lines[0]
	} else {
		return ""
	}
}

func sliceSummary(sl []string, l int) string {
	if len(sl) == 0 {
		return ""
	}

	var csl []string // check slice
	var b bytes.Buffer

	for _, s := range sl {
		lc := len(strings.Join(csl, ", ")) // (l)ength(c)heck
		switch {
		case lc <= l-10 && len(s) <= 20: //
			csl = append(csl, s)
		case lc <= l && len(s) <= 12:
			csl = append(csl, s)
		}
	}

	b.WriteString(strings.Join(csl, ", "))

	if len(sl) != len(csl) {
		b.WriteString("...")
	}

	return b.String()
}

// --> main fns

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

func (rs Repos) verifyDivs(e Emoji, f Flags, t *Timer) {

	// sort
	rs.sortByPath()

	// get all divs, remove duplicates
	var dvs []string  // divs
	var zdvs []string // zone divisions (go, main, google-apps-script etc)

	for _, r := range rs {
		dvs = append(dvs, r.DivPath)
		zdvs = append(zdvs, r.Division)
	}

	dvs = removeDuplicates(dvs)
	zdvs = removeDuplicates(zdvs)

	zds := sliceSummary(zdvs, 25) // zone division summary

	// print
	targetPrint(f, "%v  verifying divs [%v](%v)", e.FileCabinet, len(dvs), zds)

	// track created, verified and missing divs
	var cd []string // created divs
	var vd []string // verified divs
	var id []string // inaccessible divs // --> FLAG: change to unverified?

	for _, r := range rs {

		_, err := os.Stat(r.DivPath)

		// create div if missing and active run
		if os.IsNotExist(err) && isActive(f) {
			targetPrint(f, "%v creating %v", e.Folder, r.DivPath)
			os.MkdirAll(r.DivPath, 0777)
			cd = append(cd, r.DivPath)
		}

		// check div status
		info, err := os.Stat(r.DivPath)

		switch {
		case noPermission(info):
			r.markError(e, f, "fatal: No permsission", "verify-divs")
			id = append(id, r.DivPath)
		case !info.IsDir():
			r.markError(e, f, "fatal: File occupying path", "verify-divs")
			id = append(id, r.DivPath)
		case os.IsNotExist(err):
			r.markError(e, f, "fatal: No directory", "verify-divs")
			id = append(id, r.DivPath)
		case err != nil:
			r.markError(e, f, "fatal: No directory", "verify-divs")
			id = append(id, r.DivPath)
		default:
			r.Verified = true
			vd = append(vd, r.DivPath)
		}
	}

	// timer
	t.markMoment("verify-divs")

	// remove duplicates from slices
	vd = removeDuplicates(vd)
	id = removeDuplicates(id)

	// summary
	var b bytes.Buffer

	if len(dvs) == len(vd) {
		b.WriteString(e.Briefcase)
	} else {
		b.WriteString(e.Slash)
	}

	b.WriteString(" [")
	b.WriteString(strconv.Itoa(len(vd)))
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

func (rs Repos) verifyCloned(e Emoji, f Flags, t *Timer) {
	var pc []string // pending clone

	for _, r := range rs {
		r.gitCheckPending(e, f)

		if r.PendingClone == true {
			pc = append(pc, r.Name)
		}
	}

	// return if there are no pending repos

	if len(pc) <= 0 {
		return
	}

	// if there are pending repos
	targetPrint(f, "%v cloning [%v]", e.Sheep, len(pc))

	// verify each repo (async)
	var wg sync.WaitGroup
	for i := range rs {
		wg.Add(1)
		go func(r *Repo) {
			defer wg.Done()
			r.gitClone(e, f)
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
	t.markMoment("verify-repos")

	// summary
	var b bytes.Buffer

	b.WriteString(e.Truck)
	b.WriteString(" [")
	b.WriteString(strconv.Itoa(len(cr)))
	b.WriteString("/")
	b.WriteString(strconv.Itoa(len(pc)))
	b.WriteString("] cloned")

	tr := time.Millisecond // truncate
	b.WriteString(" {")
	b.WriteString(t.getSplit().Truncate(tr).String())
	b.WriteString(" / ")
	b.WriteString(t.getTime().Truncate(tr).String())
	b.WriteString("}")

	targetPrint(f, b.String())
}

func (rs Repos) verifyRepos(e Emoji, f Flags, t *Timer) {
	var rn []string // repo names

	for _, r := range rs {
		rn = append(rn, r.Name)
	}

	rns := sliceSummary(rn, 25)

	// print
	targetPrint(f, "%v  verifying repos [%v](%v)", e.Satellite, len(rs), rns)

	// verify each repo (async)
	var wg sync.WaitGroup
	for i := range rs {
		wg.Add(1)
		go func(r *Repo) {
			defer wg.Done()
			r.gitConfigOriginURL(e, f)
			r.gitRemoteUpdate(e, f)
			r.gitAbbrevRef(e, f)
			r.gitLocalSHA(e, f)
			r.gitUpstreamSHA(e, f)
			r.gitMergeBaseSHA(e, f)
			r.gitRevParseUpstream(e, f)
			r.gitDiffsNameOnly(e, f)
			r.gitShortstat(e, f)
			r.gitUntracked(e, f)
			r.setStatus(e, f)
		}(rs[i])
	}
	wg.Wait()

	// track Complete, Pending, Skipped and Scheduled
	var cr []string  // complete repos
	var pr []string  // pending repos
	var sk []string  // skipped repos
	var sch []string // scheduled repos

	for _, r := range rs {
		if r.Category == "Complete" {
			cr = append(cr, r.Name)
		}

		if r.Category == "Pending" {
			pr = append(pr, r.Name)
		}

		if r.Category == "Skipped" {
			sk = append(sk, r.Name)
		}

		if r.Category == "Scheduled" {
			sch = append(sch, r.Name)
		}
	}

	// timer
	t.markMoment("verify-repos")

	var b bytes.Buffer

	if len(cr) == len(rs) {
		b.WriteString(e.Checkmark)
	} else {
		b.WriteString(e.Traffic)
	}

	b.WriteString(" [")
	b.WriteString(strconv.Itoa(len(cr)))
	b.WriteString("/")
	b.WriteString(strconv.Itoa(len(rs)))
	b.WriteString("] complete {")

	tr := time.Millisecond // truncate
	b.WriteString(t.getSplit().Truncate(tr).String())
	b.WriteString(" / ")
	b.WriteString(t.getTime().Truncate(tr).String())
	b.WriteString("}")

	targetPrint(f, b.String())

	// scheduled repo info
	if len(sch) >= 1 {
		b.Reset()
		schs := sliceSummary(sch, 15) // scheduled repo summary
		b.WriteString(e.TimerClock)
		b.WriteString("  [")
		b.WriteString(strconv.Itoa(len(sch)))

		if loginMode(f) {
			b.WriteString("] pull scheduled (")

		} else if logoutMode(f) {
			b.WriteString("] push scheduled (")
		}

		b.WriteString(schs)
		b.WriteString(")")
		targetPrint(f, b.String())
	}

	// skipped repo info
	if len(sk) >= 1 {
		b.Reset()
		sks := sliceSummary(sk, 15) // skipped repo summary
		b.WriteString(e.Slash)
		b.WriteString(" [")
		b.WriteString(strconv.Itoa(len(sk)))
		b.WriteString("] skipped (")
		b.WriteString(sks)
		b.WriteString(")")
		targetPrint(f, b.String())
	}

	// pending repo info
	if len(pr) >= 1 {
		b.Reset()
		prs := sliceSummary(pr, 15) // pending repo summary
		b.WriteString(e.Warning)
		b.WriteString(" [")
		b.WriteString(strconv.Itoa(len(pr)))
		b.WriteString("] pending (")
		b.WriteString(prs)
		b.WriteString(")")
		targetPrint(f, b.String())
	}

}

func (rs Repos) verifyChanges(e Emoji, f Flags, t *Timer) {

	prs := initPendingRepos(rs)

	if len(prs) >= 1 {
		for _, r := range prs {

			var b bytes.Buffer

			switch r.Status {
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
				b.WriteString(" is behind ")
				b.WriteString(r.UpstreamBranch)
			case "Dirty", "DirtyUntracked", "DirtyAhead", "DirtyBehind":
				b.WriteString(e.Pig)
				b.WriteString(" ")
				b.WriteString(r.Name)
				b.WriteString(" is dirty [")
				b.WriteString(strconv.Itoa((len(r.DiffsNameOnly))))
				b.WriteString("]{")
				b.WriteString(r.DiffsSummary)
				b.WriteString("}(")
				b.WriteString(r.ShortStatSummary)
				b.WriteString(")")
			case "Untracked", "UntrackedAhead", "UntrackedBehind":
				b.WriteString(e.Pig)
				b.WriteString(" ")
				b.WriteString(r.Name)
				b.WriteString(" is untracked [")
				b.WriteString(strconv.Itoa(len(r.UntrackedFiles)))
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

			switch r.Status {
			case "DirtyUntracked":
				b.WriteString(" and untracked [")
				b.WriteString(strconv.Itoa(len(r.UntrackedFiles)))
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
				b.WriteString(" & is behind ")
				b.WriteString(r.UpstreamBranch)
			}

			targetPrint(f, b.String())

			switch r.Status {
			case "Ahead":
				fmt.Printf("%v push changes to %v? ", e.Rocket, r.Remote)
			case "Behind":
				fmt.Printf("%v pull changes from %v? ", e.Ship, r.Remote)
			case "Dirty":
				fmt.Printf("%v add all files, commit and push to %v? ", e.Clipboard, r.Remote)
			case "DirtyUntracked":
				fmt.Printf("%v add all files, commit and push to %v? ", e.Clipboard, r.Remote)
			case "DirtyAhead":
				fmt.Printf("%v add all files, commit and push to %v? ", e.Clipboard, r.Remote)
			case "DirtyBehind":
				fmt.Printf("%v stash all files, pull changes, commit and push to %v? ", e.Clipboard, r.Remote)
			case "Untracked":
				fmt.Printf("%v add all files, commit and push to %v? ", e.Clipboard, r.Remote)
			case "UntrackedAhead":
				fmt.Printf("%v add all files, commit and push to %v? ", e.Clipboard, r.Remote)
			case "UntrackedBehind":
				fmt.Printf("%v stash all files, pull changes, commit and push to %v? ", e.Clipboard, r.Remote)
			}

			// prompt for approval
			r.checkConfirmed()

			// prompt for commit message
			if r.Category != "Skipped" && strings.Contains(r.GitAction, "commit") {
				fmt.Printf("%v commit message: ", e.Memo)
				r.checkCommitMessage()
			}
		}
	}
}

func (rs Repos) submitChanges(e Emoji, f Flags, t *Timer) {
	srs := initScheludedRepos(rs)
	// fmt.Printf("len %v\n", len(srs))

	if len(srs) >= 1 {
		for _, r := range srs {
			fmt.Println(r.Name)
		}
	}

}

// debugger
func (rs Repos) debug() {
	for _, r := range rs {
		if r.ErrorShort != "" {
			fmt.Printf("%v|%v (%v)\n", r.Name, r.ErrorName, r.ErrorFirst)
			fmt.Printf("%v\n", r.ErrorShort)
		}
	}
}

func main() {
	e, f, rs, t := initRun()
	rs.verifyDivs(e, f, t)
	rs.verifyCloned(e, f, t)
	rs.verifyRepos(e, f, t)
	rs.verifyChanges(e, f, t)
	rs.submitChanges(e, f, t)
}
