package main_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/wchan2/bloodhound"
)

var _ = Describe(`HTTPTrafficFilter`, func() {
	Describe(`#Filter`, func() {
		var (
			trafficFilter TrafficFilter

			event    Event
			filtered bool
		)

		BeforeEach(func() {
			trafficFilter = HTTPTrafficFilter
		})

		Context(`when a packet is a http request`, func() {
			JustBeforeEach(func() {
				packet := new(packetMock)
				appLayer := new(appLayerMock)
				appLayer.payload = []byte(sampleRequest)
				packet.appLayer = appLayer
				event, filtered = trafficFilter.Filter(packet)
			})

			It(`returns true`, func() {
				Expect(filtered).To(BeTrue())
			})

			It(`returns a populated event with the current time`, func() {
				Expect(event.Time).To(BeTemporally("<=", time.Now()))
			})

			It(`returns a populated event with the correct destination`, func() {
				Expect(event.Destination).To(Equal("a.disquscdn.com/next/embed/alfie.f51946af45e0b561c60f768335c9eb79.js"))
			})
		})

		Context(`when a packet is not a http request`, func() {
			JustBeforeEach(func() {
				packet := new(packetMock)
				appLayer := new(appLayerMock)
				appLayer.payload = []byte(`non-http request`)
				packet.appLayer = appLayer
				event, filtered = trafficFilter.Filter(packet)
			})

			It(`returns an empty event`, func() {
				Expect(event).To(Equal(Event{}))
			})

			It(`returns false`, func() {
				Expect(filtered).To(BeFalse())
			})
		})
	})
})
