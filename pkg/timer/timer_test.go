package timer

import (
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	tt := Init()        // test Timer
	sm := tt.Moments[0] // Timer initializes with a (S)tart (M)oment

	if sm.Name != "Start" {
		t.Errorf("TestTimer: Start Moment Name error")
	}

	n := time.Now() // now

	if sm.Time.UnixNano() > n.UnixNano() {
		t.Errorf("TestTimer: Start Time error")
	}

	for _, c := range []struct {
		name string
	}{
		{"TestMoment1"},
		{"TestMoment2"},
		{"TestMoment3"},
	} {
		tt.Mark(c.name)
	}

	// Timer.Moments should have 4 Moments: Start, TestMoment1, ...

	if len(tt.Moments) != 4 {
		t.Errorf("TestTimer: Moments length error")
	}

	for _, c := range []struct {
		name string
	}{
		{"TestMoment1"},
		{"TestMoment2"},
		{"TestMoment3"},
	} {
		if m, err := tt.GetMoment(c.name); err != nil {
			t.Errorf("TestTimer: GetMoment error (%v)", c.name)
		} else {
			switch {
			case m.Start < 0:
				t.Errorf("TestTimer: Start error (%v)", c.name)
			case m.Split < 0:
				t.Errorf("TestTimer: Split error (%v)", c.name)
			}
		}
	}

	if _, err := tt.GetMoment("UndefinedMoment4"); err == nil {
		t.Errorf("TestTimer: GetMoment didn't error w/ UndefinedMoment4")
	}

}
