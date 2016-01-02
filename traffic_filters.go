package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/google/gopacket"
)

var HTTPTrafficFilter = NewTrafficFilter(func(packet gopacket.Packet) (Event, bool) {
	appLayer := packet.ApplicationLayer()
	if appLayer != nil && strings.Contains(string(appLayer.Payload()), "HTTP") {
		payloadReader := bytes.NewReader(appLayer.Payload())
		bufferedPayloadReader := bufio.NewReader(payloadReader)
		request, err := http.ReadRequest(bufferedPayloadReader)
		if err != nil && err != io.EOF {
			return Event{}, false
		}
		payload, err := ioutil.ReadAll(request.Body)
		if err != nil {
			return Event{}, false
		}
		return Event{
			Payload:     payload,
			Destination: fmt.Sprintf("%s%s", request.Host, request.RequestURI),
			Time:        time.Now(),
		}, true
	}
	return Event{}, false
})

type TrafficFilter interface {
	Filter(gopacket.Packet) (Event, bool)
}

type PacketFilterFn func(gopacket.Packet) (Event, bool)

type trafficFilter struct {
	packetFilterFn PacketFilterFn
}

func NewTrafficFilter(packetFilterFn PacketFilterFn) *trafficFilter {
	return &trafficFilter{packetFilterFn: packetFilterFn}
}

func (t *trafficFilter) Filter(packet gopacket.Packet) (Event, bool) {
	return t.packetFilterFn(packet)
}
