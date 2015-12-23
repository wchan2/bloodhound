package main

import "log"

type ConsoleAlert struct {
	threshold      Threshold
	events         []Event
	alertTriggered bool
}

func NewConsoleAlert(threshold Threshold) *ConsoleAlert {
	return &ConsoleAlert{threshold: threshold, events: []Event{}}
}

func (c *ConsoleAlert) Check(event Event) {
	c.events = append(c.events, event)
	c.prune(event)
	if len(c.events) >= c.threshold.Hits && !c.alertTriggered {
		log.Printf("High traffic generated an alert - hits = 2, triggered at %s", event.Time.String())
		c.alertTriggered = true
	} else if len(c.events) < c.threshold.Hits && c.alertTriggered {
		log.Printf("Traffic returned to normal, triggered at %s\n", event.Time.String())
		c.alertTriggered = false
	}
}

func (c *ConsoleAlert) prune(event Event) {
	var lastGreatestIndex int
	for i := len(c.events) - 1; i >= 0; i-- {
		if event.Time.Sub(c.events[i].Time) > c.threshold.Duration {
			lastGreatestIndex = i + 1
			break
		}
	}
	c.events = c.events[lastGreatestIndex:]
}
