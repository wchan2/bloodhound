package main

import (
	"bufio"
	"bytes"
	"fmt"
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
		req, err := http.ReadRequest(bufferedPayloadReader)
		fmt.Println(string(appLayer.Payload()))
		if err != nil && err != io.EOF {
			fmt.Println(err.Error())
			fmt.Println("HERE")
			return false
		}
		if err == io.EOF {
			fmt.Println("REQUEST URI", req.RequestURI)
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
