package main

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/google/gopacket"
)

var HTTPTrafficFilter = NewTrafficFilter(func(packet gopacket.Packet) bool {
	appLayer := packet.ApplicationLayer()
	if appLayer != nil && strings.Contains(string(appLayer.Payload()), "HTTP") {
		payloadReader := bytes.NewReader(appLayer.Payload())
		bufferedPayloadReader := bufio.NewReader(payloadReader)
		_, err := http.ReadRequest(bufferedPayloadReader)
		if err != nil && err != io.EOF {
			return false
		}
		return true
	}
	return false
})

type TrafficFilter interface {
	Filter(gopacket.Packet) bool
}

type PacketFilterFn func(gopacket.Packet) bool

type trafficFilter struct {
	packetFilterFn PacketFilterFn
}

func NewTrafficFilter(packetFilterFn PacketFilterFn) *trafficFilter {
	return &trafficFilter{packetFilterFn: packetFilterFn}
}

func (t *trafficFilter) Filter(packet gopacket.Packet) bool {
	return t.packetFilterFn(packet)
}
