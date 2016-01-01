package main

import (
	"fmt"
	"time"
)

type Alert interface {
	Check(Event)
}

type TotalTrafficAlert struct {
	hits         int
	duration     time.Duration
	notification Notification

	alertTriggered bool
	events         []Event
}

func NewTotalTrafficAlert(hits int, duration time.Duration, notification Notification) *TotalTrafficAlert {
	return &TotalTrafficAlert{
		hits:         hits,
		duration:     duration,
		notification: notification,
	}
}

func (t *TotalTrafficAlert) Check(event Event) {
	t.add(event)

	if t.isExceeded() && !t.alertTriggered {
		t.notification.Send(fmt.Sprintf("High traffic generated an alert - hits = %d, triggered at %s", len(t.events), event.Time.String()))
		t.alertTriggered = true
	} else if !t.isExceeded() && t.alertTriggered {
		reason := fmt.Sprintf("Traffic returned to normal, triggered at %s", event.Time.String())
		t.notification.Send(reason)
		t.alertTriggered = false
	}
}

func (t *TotalTrafficAlert) add(event Event) {
	t.events = append(t.events, event)
	t.pruneUpTo(event.Time)
}

func (t *TotalTrafficAlert) isExceeded() bool {
	return len(t.events) > 0 && len(t.events) >= t.hits
}

func (t *TotalTrafficAlert) pruneUpTo(currentTime time.Time) {
	var lastGreatestIndex int
	for i := len(t.events) - 1; i >= 0; i-- {
		if currentTime.Sub(t.events[i].Time) > t.duration {
			lastGreatestIndex = i + 1
			break
		}
	}
	t.events = t.events[lastGreatestIndex:]
}
