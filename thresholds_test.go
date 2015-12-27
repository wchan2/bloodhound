package main_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/wchan2/bloodhound"
)

type alertMock struct {
	message string
}

func (s *alertMock) Send(message string) {
	s.message = message
}

var _ = Describe(`TotalTrafficThreshold`, func() {
	var (
		threshold Threshold
		alert     *alertMock
	)

	JustBeforeEach(func() {
		threshold = NewTotalTrafficThreshold(2, 2*time.Minute, alert)
	})

	Describe(`#Check`, func() {
		BeforeEach(func() {
			alert = new(alertMock)
		})

		Context(`when the threshold is exceeded`, func() {
			var currentTime, twoMinutesAgo time.Time

			BeforeEach(func() {
				currentTime = time.Now()
				twoMinutesAgo = currentTime.Add(-2 * time.Minute)
			})

			JustBeforeEach(func() {
				threshold.Check(Event{Time: twoMinutesAgo})
				threshold.Check(Event{Time: currentTime})
			})

			It(`sends an alert message`, func() {
				alertMessage := fmt.Sprintf("High traffic generated an alert - hits = 2, triggered at %s", currentTime.String())
				Expect(alert.message).To(Equal(alertMessage))
			})
		})

		Context(`when the threshold is not exceeded`, func() {
			var currentTime, fourMinutesAgo time.Time

			BeforeEach(func() {
				currentTime = time.Now()
				fourMinutesAgo = currentTime.Add(-4 * time.Minute)
			})

			JustBeforeEach(func() {
				threshold.Check(Event{Time: fourMinutesAgo})
				threshold.Check(Event{Time: currentTime})
			})

			It(`does not send an alert message`, func() {
				Expect(alert.message).To(BeEmpty())
			})
		})

		Context(`when the threshold is resolved`, func() {
			var currentTime, threeMinutesAgo, fourMinutesAgo time.Time

			BeforeEach(func() {
				currentTime = time.Now()
				threeMinutesAgo = currentTime.Add(-3 * time.Minute)
				fourMinutesAgo = currentTime.Add(-4 * time.Minute)
			})

			JustBeforeEach(func() {
				// trigger an alert
				threshold.Check(Event{Time: fourMinutesAgo})
				threshold.Check(Event{Time: threeMinutesAgo})

				// resolve the alert
				threshold.Check(Event{Time: currentTime})
			})

			It(`sends a resolved message`, func() {
				resolvedMessage := fmt.Sprintf("Traffic returned to normal, triggered at %s", currentTime.String())
				Expect(alert.message).To(Equal(resolvedMessage))
			})
		})
	})
})
