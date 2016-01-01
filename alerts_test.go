package main_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/wchan2/bloodhound"
)

type notificationMock struct {
	message string
}

func (s *notificationMock) Send(message string) {
	s.message = message
}

var _ = Describe(`TotalTrafficAlert`, func() {
	var (
		alert        Alert
		notification *notificationMock
	)

	JustBeforeEach(func() {
		alert = NewTotalTrafficAlert(2, 2*time.Minute, notification)
	})

	Describe(`#Check`, func() {
		BeforeEach(func() {
			notification = new(notificationMock)
		})

		Context(`when the threshold is exceeded`, func() {
			var currentTime, twoMinutesAgo time.Time

			BeforeEach(func() {
				currentTime = time.Now()
				twoMinutesAgo = currentTime.Add(-2 * time.Minute)
			})

			JustBeforeEach(func() {
				alert.Check(Event{Time: twoMinutesAgo})
				alert.Check(Event{Time: currentTime})
			})

			It(`sends an alert message`, func() {
				notificationMessage := fmt.Sprintf("High traffic generated an alert - hits = 2, triggered at %s", currentTime.String())
				Expect(notification.message).To(Equal(notificationMessage))
			})
		})

		Context(`when the threshold is not exceeded`, func() {
			var currentTime, fourMinutesAgo time.Time

			BeforeEach(func() {
				currentTime = time.Now()
				fourMinutesAgo = currentTime.Add(-4 * time.Minute)
			})

			JustBeforeEach(func() {
				alert.Check(Event{Time: fourMinutesAgo})
				alert.Check(Event{Time: currentTime})
			})

			It(`does not send an alert message`, func() {
				Expect(notification.message).To(BeEmpty())
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
				alert.Check(Event{Time: fourMinutesAgo})
				alert.Check(Event{Time: threeMinutesAgo})

				// resolve the alert
				alert.Check(Event{Time: currentTime})
			})

			It(`sends a resolved message`, func() {
				resolvedMessage := fmt.Sprintf("Traffic returned to normal, triggered at %s", currentTime.String())
				Expect(notification.message).To(Equal(resolvedMessage))
			})
		})
	})
})
