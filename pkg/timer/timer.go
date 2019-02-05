// Package timer records Moments.
package timer

import (
	"errors"
	"time"
)

// Moment marks moments in time.
type Moment struct {
	Name  string
	Time  time.Time
	Start time.Duration // duration since start
	Split time.Duration // duration since last moment
}

// Timer collects Moments.
type Timer struct {
	Moments []Moment
}

// Init initializes a *Timer with a Start Moment.
func Init() *Timer {
	ti := new(Timer)
	st := Moment{Name: "Start", Time: time.Now()} // (st)art
	ti.Moments = append(ti.Moments, st)
	return ti
}

// Mark marks a moment in time as a Moment and appends t.Moments.
func (ti *Timer) Mark(s string) {
	sm := ti.Moments[0]                          // (s)tarting (m)oment
	lm := ti.Moments[len(ti.Moments)-1]          // (l)ast (m)oment
	m := Moment{Name: s, Time: time.Now()}       // name and time
	m.Start = time.Since(sm.Time).Truncate(1000) // duration since start
	m.Split = m.Start - lm.Start                 // duration since last moment
	ti.Moments = append(ti.Moments, m)           // append Moment
}

// Time returns the elapsed time at the last recorded moment in *Timer.
func (ti *Timer) Time() time.Duration {
	lm := ti.Moments[len(ti.Moments)-1] // (l)ast (m)oment
	return lm.Start
}

// Split returns the split time for the last recorded moment in *Timer.
func (ti *Timer) Split() time.Duration {
	lm := ti.Moments[len(ti.Moments)-1] // (l)ast (m)oment
	return lm.Split
}

// Get returns a Moment and an error value from a *Timer.
func (ti *Timer) Get(s string) (Moment, error) {
	for _, m := range ti.Moments {
		if m.Name == s {
			return m, nil
		}
	}

	var em Moment // (e)mpty (m)oment
	return em, errors.New("no moment found")
}
