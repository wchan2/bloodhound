package main

import (
	"fmt"
	"strings"
	"time"
)

type TrafficMonitor interface {
	Monitor(Event)
	Stop()
}

type TrafficStatistics struct {
	Destination    string
	AveragePayload float64
	TotalPayload   int64
	Count          int64
}

func (t *TrafficStatistics) String() string {
	statisticsFormat := "Destination: %s\nAverage Payload: %f\nTotal Payload: %d\n"
	return fmt.Sprintf(statisticsFormat, t.Destination, t.AveragePayload, t.TotalPayload)
}

type SummaryStatsTrafficMonitor struct {
	duration     time.Duration
	notification Notification

	ticker *time.Ticker
	events chan Event

	statistics map[string]TrafficStatistics
}

func NewSummaryStatsTrafficMonitor(duration time.Duration, notification Notification) *SummaryStatsTrafficMonitor {
	monitor := &SummaryStatsTrafficMonitor{
		duration:     duration,
		notification: notification,
		events:       make(chan Event),

		statistics: map[string]TrafficStatistics{},
	}
	go monitor.consumeEvents()
	go monitor.publishStatistics()
	return monitor
}

func (s *SummaryStatsTrafficMonitor) Monitor(event Event) {
	s.events <- event
}

func (s *SummaryStatsTrafficMonitor) Stop() {
	s.ticker.Stop()
	close(s.events)
}

func (s *SummaryStatsTrafficMonitor) summary() string {
	statistics := make([]string, len(s.statistics))
	i := 0
	for _, statistic := range s.statistics {
		statistics[i] = statistic.String()
		i++
	}
	return strings.Join(statistics, "")
}

func (s *SummaryStatsTrafficMonitor) publishStatistics() {
	s.ticker = time.NewTicker(s.duration)
	for _ = range s.ticker.C {
		s.notification.Send(s.summary())
		s.statistics = map[string]TrafficStatistics{}
	}
}

func (s *SummaryStatsTrafficMonitor) consumeEvents() {
	for event := range s.events {
		if statistic, ok := s.statistics[event.Destination]; ok {
			s.statistics[event.Destination] = TrafficStatistics{
				Destination:    event.Destination,
				AveragePayload: (statistic.AveragePayload + float64(len(event.Payload))) / float64((statistic.Count + 1)),
				TotalPayload:   statistic.TotalPayload + int64(len(event.Payload)),
				Count:          statistic.Count + 1,
			}
		} else {
			s.statistics[event.Destination] = TrafficStatistics{
				Destination:    event.Destination,
				AveragePayload: float64(len(event.Payload)),
				TotalPayload:   int64(len(event.Payload)),
				Count:          1,
			}
		}
	}
}
