package timer

import (
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	ti := Init()
	sm := ti.Moments[0]

	if sm.Name != "Start" {
		t.Errorf("TestTimer: Start Moment Name error")
	}

	n := time.Now() // now

	if sm.Time.UnixNano() > n.UnixNano() {
		t.Errorf("TestTimer: Start Time error")
	}

	for _, tr := range []struct {
		name string
	}{
		{"TestMoment1"},
		{"TestMoment2"},
		{"TestMoment3"},
	} {
		ti.Mark(tr.name)
	}

	// Timer.Moments should have 4 Moments: Start, TestMoment1, ...

	if len(ti.Moments) != 4 {
		t.Errorf("TestTimer: Moments length error")
	}

	for _, tr := range []struct {
		name string
	}{
		{"TestMoment1"},
		{"TestMoment2"},
		{"TestMoment3"},
	} {
		if m, err := ti.Get(tr.name); err != nil {
			t.Errorf("TestTimer: GetMoment error (%v)", tr.name)
		} else {
			switch {
			case m.Start < 0:
				t.Errorf("TestTimer: Start error (%v)", tr.name)
			case m.Split < 0:
				t.Errorf("TestTimer: Split error (%v)", tr.name)
			}
		}
	}

	if _, err := ti.Get("UndefinedMoment4"); err == nil {
		t.Errorf("TestTimer: GetMoment didn't error w/ UndefinedMoment4")
	}

}
