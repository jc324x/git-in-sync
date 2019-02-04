// Package e adds emoji to git-in-sync.
package e

import (
	"fmt"
	"html"
	"os"
	"os/exec"
	"strconv"
)

var em = map[string]int{
	"AlarmClock":           9200,
	"Boat":                 128676,
	"Book":                 128214,
	"Books":                128218,
	"Box":                  128230,
	"Briefcase":            128188,
	"BuildingConstruction": 127959,
	"Bunny":                128048,
	"Checkmark":            9989,
	"Clapper":              127916,
	"Clipboard":            128203,
	"CrystalBall":          128302,
	"DirectHit":            127919,
	"Desert":               127964,
	"FaxMachine":           128224,
	"Finger":               128073,
	"FileCabinet":          128452,
	"Flag":                 127937,
	"FlagInHole":           9971,
	"Fire":                 128293,
	"Folder":               128193,
	"Glasses":              128083,
	"Hole":                 128371,
	"Hourglass":            9203,
	"Inbox":                128229,
	"Microscope":           128300,
	"Memo":                 128221,
	"Outbox":               128228,
	"Pager":                128223,
	"Parents":              128106,
	"Pen":                  128394,
	"Pig":                  128055,
	"Popcorn":              127871,
	"Rocket":               128640,
	"Run":                  127939,
	"Satellite":            128752,
	"SatelliteDish":        128225,
	"Slash":                128683,
	"Ship":                 128674,
	"Sheep":                128017,
	"Squirrel":             128063,
	"Telescope":            128301,
	"Text":                 128172,
	"ThumbsUp":             128077,
	"TimerClock":           9202,
	"Traffic":              128678,
	"Truck":                128666,
	"Turtle":               128034,
	"Unicorn":              129412,
	"Warning":              128679,
}

// convert returns an emoji character as a string value.
func convert(n int) string {
	return html.UnescapeString("&#" + strconv.Itoa(n) + ";")
}

// Get returns an emoji character as a string.
func Get(s string) string {

	if val, ok := em[s]; ok {
		return convert(val)
	}

	return "#"
}

// Count returns the number of emoji in the map.
func Count() int {
	return len(em)
}

// ClearScreen clears the screen.
func ClearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// printAll prints values in map em.
func printAll() {
	for _, val := range em {
		fmt.Printf("%v: \"%v\"\n", val, convert(val))
	}
}
