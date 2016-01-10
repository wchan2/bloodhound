package main_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/wchan2/bloodhound"
)

var _ = Describe(`SummaryStatsTrafficMonitor`, func() {
	Describe(`#Monitor`, func() {
		var (
			trafficMonitor TrafficMonitor
			notification   *notificationMock

			events []Event
		)

		BeforeEach(func() {
			notification = new(notificationMock)
			trafficMonitor = NewSummaryStatsTrafficMonitor(1*time.Second, notification)
			currentTime := time.Now()
			events = []Event{
				{
					Destination: `destination1`,
					Payload:     []byte(`random payload`),
					Time:        currentTime,
				},
				{
					Destination: `destination1`,
					Payload:     []byte(`another random payload`),
					Time:        currentTime.Add(1 * time.Second),
				},
				{
					Destination: `destination2`,
					Payload:     []byte(`yet another random payload`),
					Time:        currentTime.Add(2 * time.Second),
				},
			}
		})

		JustBeforeEach(func() {
			for _, event := range events {
				trafficMonitor.Monitor(event)
			}
		})

		It(`calculates the summary statistics`, func() {
			destination1TotalPayload := len(events[0].Payload) + len(events[1].Payload)
			destination1AvgPayload := float64(destination1TotalPayload) / 2

			destination2TotalPayload := len(events[2].Payload)
			destination2AvgPayload := float64(destination2TotalPayload) / 1

			statisticsFormat := "Destination: %s\nAverage Payload: %f\nTotal Payload: %d\nCount: %d\n"
			Eventually(func() string { return notification.message }, 2*time.Second, 500*time.Millisecond).Should(And(
				ContainSubstring(fmt.Sprintf(statisticsFormat, "destination1", destination1AvgPayload, destination1TotalPayload, 2)),
				ContainSubstring(fmt.Sprintf(statisticsFormat, "destination2", destination2AvgPayload, destination2TotalPayload, 1)),
			))
		})
	})
})
