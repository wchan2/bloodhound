package main_test

import (
	"github.com/google/gopacket"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/wchan2/bloodhound"
)

var sampleRequest = `GET /next/embed/alfie.f51946af45e0b561c60f768335c9eb79.js HTTP/1.1
Host: a.disquscdn.com
Connection: keep-alive
Accept: */*
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/46.0.2490.80 Safari/537.36
DNT: 1
Referer: http://www.lucainvernizzi.net/blog/2015/02/12/extracting-urls-from-network-traffic-in-just-9-python-lines-with-scapy-http/
Accept-Encoding: gzip, deflate, sdch
Accept-Language: en-US,en;q=0.8

`

type packetMock struct {
	appLayer gopacket.ApplicationLayer
}

func (p *packetMock) ApplicationLayer() gopacket.ApplicationLayer { return p.appLayer }

func (p *packetMock) String() string                                  { return "" }
func (p *packetMock) Dump() string                                    { return "" }
func (p *packetMock) Layers() []gopacket.Layer                        { return nil }
func (p *packetMock) Layer(_ gopacket.LayerType) gopacket.Layer       { return nil }
func (p *packetMock) LayerClass(_ gopacket.LayerClass) gopacket.Layer { return nil }
func (p *packetMock) LinkLayer() gopacket.LinkLayer                   { return nil }
func (p *packetMock) NetworkLayer() gopacket.NetworkLayer             { return nil }
func (p *packetMock) TransportLayer() gopacket.TransportLayer         { return nil }
func (p *packetMock) ErrorLayer() gopacket.ErrorLayer                 { return nil }
func (p *packetMock) Data() []byte                                    { return nil }
func (p *packetMock) Metadata() *gopacket.PacketMetadata              { return nil }

type appLayerMock struct {
	payload []byte
}

func (a *appLayerMock) LayerType() gopacket.LayerType { return 0 }
func (a *appLayerMock) LayerContents() []byte         { return nil }
func (a *appLayerMock) LayerPayload() []byte          { return nil }
func (a *appLayerMock) Payload() []byte               { return a.payload }

var _ = Describe(`HTTPTrafficFilter`, func() {
	Describe(`#Filter`, func() {
		var (
			trafficFilter TrafficFilter
			filtered      bool
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
				filtered = trafficFilter.Filter(packet)
			})

			It(`returns true`, func() {
				Expect(filtered).To(BeTrue())
			})
		})

		Context(`when a packet is not a http request`, func() {
			JustBeforeEach(func() {
				packet := new(packetMock)
				appLayer := new(appLayerMock)
				appLayer.payload = []byte(`non-http request`)
				packet.appLayer = appLayer
				filtered = trafficFilter.Filter(packet)
			})

			It(`returns false`, func() {
				Expect(filtered).To(BeFalse())
			})
		})
	})
})
