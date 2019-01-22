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
	t := new(Timer)
	st := Moment{Name: "Start", Time: time.Now()} // (st)art
	t.Moments = append(t.Moments, st)
	return t
}

// MarkMoment marks a moment in time as a Moment and appends t.Moments.
func (t *Timer) MarkMoment(s string) {
	sm := t.Moments[0]                           // (s)tarting (m)oment
	lm := t.Moments[len(t.Moments)-1]            // (l)ast (m)oment
	m := Moment{Name: s, Time: time.Now()}       // name and time
	m.Start = time.Since(sm.Time).Truncate(1000) // duration since start
	m.Split = m.Start - lm.Start                 // duration since last moment
	t.Moments = append(t.Moments, m)             // append Moment
}

// GetTime returns the elapsed time at the last recorded moment in *Timer.
func (t *Timer) Time() time.Duration {
	lm := t.Moments[len(t.Moments)-1] // (l)ast (m)oment
	return lm.Start
}

// GetSplit returns the split time for the last recorded moment in *Timer.
func (t *Timer) Split() time.Duration {
	lm := t.Moments[len(t.Moments)-1] // (l)ast (m)oment
	return lm.Split
}

// GetMoment returns a Moment and an error value from a *Timer.
func (t *Timer) GetMoment(s string) (Moment, error) {
	for _, m := range t.Moments {
		if m.Name == s {
			return m, nil
		}
	}

	var em Moment // (e)mpty (m)oment
	return em, errors.New("no moment found")
}
