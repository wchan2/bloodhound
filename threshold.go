package main

import "time"

type Threshold interface {
	Add(Event)
	IsExceeded() bool
	IsResolved() bool
}

type TotalTrafficThreshold struct {
	hits     int
	duration time.Duration
	events   []Event
}

func NewTotalTrafficThreshold(hits int, duration time.Duration) *TotalTrafficThreshold {
	return &TotalTrafficThreshold{
		hits:     hits,
		duration: duration,
		events:   []Event{},
	}
}

func (t *TotalTrafficThreshold) Add(event Event) {
	t.events = append(t.events, event)
	t.pruneUpTo(event.Time)
}

func (t *TotalTrafficThreshold) IsExceeded() bool {
	return len(t.events) >= t.hits
}

func (t *TotalTrafficThreshold) IsResolved() bool {
	return len(t.events) < t.hits
}

func (t *TotalTrafficThreshold) pruneUpTo(currentTime time.Time) {
	var lastGreatestIndex int
	for i := len(t.events) - 1; i >= 0; i-- {
		if currentTime.Sub(t.events[i].Time) > t.duration {
			lastGreatestIndex = i + 1
			break
		}
	}
	t.events = t.events[lastGreatestIndex:]
}
