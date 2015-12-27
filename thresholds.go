package main

import (
	"fmt"
	"time"
)

type Threshold interface {
	Check(Event)
}

type TotalTrafficThreshold struct {
	hits     int
	duration time.Duration
	alert    Alert

	alertTriggered bool
	events         []Event
}

func NewTotalTrafficThreshold(hits int, duration time.Duration, alert Alert) *TotalTrafficThreshold {
	return &TotalTrafficThreshold{
		hits:     hits,
		duration: duration,
		alert:    alert,
	}
}

func (t *TotalTrafficThreshold) Check(event Event) {
	t.add(event)

	if t.isExceeded() && !t.alertTriggered {
		t.alert.Send(fmt.Sprintf("High traffic generated an alert - hits = %d, triggered at %s", len(t.events), event.Time.String()))
		t.alertTriggered = true
	} else if !t.isExceeded() && t.alertTriggered {
		reason := fmt.Sprintf("Traffic returned to normal, triggered at %s", event.Time.String())
		t.alert.Send(reason)
		t.alertTriggered = false
	}
}

func (t *TotalTrafficThreshold) add(event Event) {
	t.events = append(t.events, event)
	t.pruneUpTo(event.Time)
}

func (t *TotalTrafficThreshold) isExceeded() bool {
	return len(t.events) > 0 && len(t.events) >= t.hits
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
