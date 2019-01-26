package emoji

import (
	"testing"
)

func TestPrintEmoji(t *testing.T) {
	e := Init()

	for _, c := range []struct {
		j    int
		str  string
		want string
	}{
		{9200, "AlarmClock", "⏰"},
	} {
		got := e.c.str
		if got != c.want {

		}
		// if got := AbsUser(c.in); err != nil || got != c.want {
		// 	t.Errorf("AbsUser: (%v != %v)", got, c.want)
		// }
	}

}
