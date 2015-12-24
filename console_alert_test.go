package main_test

import (
	"bytes"
	"fmt"
	"log"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/wchan2/bloodhound"
)

var _ = Describe(`ConsoleAlert`, func() {
	var (
		consoleAlert *ConsoleAlert
		buf          *bytes.Buffer
	)

	BeforeEach(func() {
		buf = &bytes.Buffer{}
		log.SetOutput(buf)
		log.SetFlags(0)

		consoleAlert = NewConsoleAlert(NewTotalTrafficThreshold(2, 1*time.Minute))
	})

	Context(`when there is a single event that does not exceed the threshold`, func() {
		var currentTime time.Time

		BeforeEach(func() {
			currentTime = time.Now()
			consoleAlert.Check(Event{Time: currentTime})
		})

		It(`does not alert`, func() {
			Expect(buf.String()).To(Equal(""))
		})
	})

	Context(`when there are consecutive events that does not exceed the threshold`, func() {
		var currentTime, twoMinutesAgo, fourMinutesAgo time.Time

		BeforeEach(func() {
			currentTime = time.Now()
			fourMinutesAgo = currentTime.Add(-4 * time.Minute)
			twoMinutesAgo = currentTime.Add(-2 * time.Minute)
			buf = &bytes.Buffer{}
			log.SetOutput(buf)
			log.SetFlags(0)

			consoleAlert.Check(Event{Time: fourMinutesAgo})
			consoleAlert.Check(Event{Time: twoMinutesAgo})
			consoleAlert.Check(Event{Time: currentTime})
		})

		It(`does not alert`, func() {
			Expect(buf.String()).To(Equal(""))
		})
	})

	Context(`when consecutive event exceeds the threshold`, func() {
		var currentTime, thirtySecondsAgo, fortyFiveSecondsAgo time.Time

		BeforeEach(func() {
			currentTime = time.Now()
			fortyFiveSecondsAgo = currentTime.Add(-45 * time.Second)
			thirtySecondsAgo = currentTime.Add(-30 * time.Second)

			consoleAlert.Check(Event{Time: fortyFiveSecondsAgo})
			consoleAlert.Check(Event{Time: thirtySecondsAgo})
		})

		It(`alerts`, func() {
			alertText := fmt.Sprintf("High traffic generated an alert - hits = 2, triggered at %s\n", thirtySecondsAgo.String())
			Expect(buf.String()).To(Equal(alertText))
		})

		It(`alerts only once when an alert has already been triggered`, func() {
			consoleAlert.Check(Event{Time: currentTime})
			alertText := fmt.Sprintf("High traffic generated an alert - hits = 2, triggered at %s\n", thirtySecondsAgo.String())
			Expect(buf.String()).To(Equal(alertText))
		})
	})

	Context(`when the consecutive event exceeds the threshold and reverts to normal`, func() {
		var currentTime, threeMinutesAgo, fourMinutesAgo time.Time

		BeforeEach(func() {
			currentTime = time.Now()
			threeMinutesAgo = currentTime.Add(-3 * time.Minute)
			fourMinutesAgo = currentTime.Add(-4 * time.Minute)

			consoleAlert.Check(Event{Time: fourMinutesAgo})
			consoleAlert.Check(Event{Time: threeMinutesAgo})
			consoleAlert.Check(Event{Time: currentTime})
		})

		It(`alerts and displays a message that traffic has reverted to normal`, func() {
			alertText := fmt.Sprintf(fmt.Sprintf("High traffic generated an alert - hits = 2, triggered at %s\n", threeMinutesAgo.String()))
			trafficNeutralizedText := fmt.Sprintf(fmt.Sprintf("Traffic returned to normal, triggered at %s\n", currentTime.String()))
			message := fmt.Sprintf("%s%s", alertText, trafficNeutralizedText)
			Expect(buf.String()).To(Equal(message))
		})
	})
})
