package main

import "log"

type ConsoleAlert struct {
	threshold      Threshold
	alertTriggered bool
}

func NewConsoleAlert(threshold Threshold) *ConsoleAlert {
	return &ConsoleAlert{threshold: threshold}
}

func (c *ConsoleAlert) Check(event Event) {
	c.threshold.Add(event)
	if c.threshold.IsExceeded() && !c.alertTriggered {
		log.Printf("High traffic generated an alert - hits = 2, triggered at %s", event.Time.String())
		c.alertTriggered = true
	} else if c.threshold.IsResolved() && c.alertTriggered {
		log.Printf("Traffic returned to normal, triggered at %s\n", event.Time.String())
		c.alertTriggered = false
	}
}
