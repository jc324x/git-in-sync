package emoji

import (
	"html"
	"os"
	"os/exec"
	"reflect"
	"strconv"
)

// printEmoji returns an emoji character as a string value.
func printEmoji(n int) string {
	return html.UnescapeString("&#" + strconv.Itoa(n) + ";")
}

// Emoji collects emoji values as strings.
// Yes, I could use a map here, but I don't care.
type Emoji struct {
	AlarmClock           string
	Boat                 string
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

// Init returns an Emoji struct with all values initialized.
func Init() (e Emoji) {
	e.AlarmClock = printEmoji(9200)
	e.Boat = printEmoji(128676)
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
	return e
}

// ClearScreen clears the screen.
func ClearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
