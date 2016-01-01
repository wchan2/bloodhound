package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

var (
	iface    = flag.String("interface", "en0", "Network interface to listen for packets")
	protocol = flag.String("protocol", "tcp", "Protocol to listen for packets")
	port     = flag.String("port", "80", "Port to listen to for packets")

	duration = flag.Int("duration", 30, "Duration in seconds that")
	traffic  = flag.Int("traffic", 100, "Traffic amount that should trigger an alert")
)

type Application struct {
	applicationConfig Config
	trafficFilter     TrafficFilter
	alert             Alert
}

type Config struct {
	NetworkInterface string
	Protocol         string
	Port             string
}

func NewApplication(config Config, filter TrafficFilter, alert Alert) *Application {
	return &Application{
		applicationConfig: config,
		trafficFilter:     filter,
		alert:             alert,
	}
}

func (a *Application) Run() {
	handle, err := pcap.OpenLive(a.applicationConfig.NetworkInterface, 1024, false, 1*time.Second)
	if err != nil {
		log.Fatalf("Unable to ", err.Error())
	}
	defer handle.Close()
	handle.SetBPFFilter(fmt.Sprintf("%s port %s", a.applicationConfig.Protocol, a.applicationConfig.Port))
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		if !a.trafficFilter.Filter(packet) {
			a.alert.Check(Event{Time: time.Now()})
		}
	}
}

func main() {
	config := Config{
		NetworkInterface: *iface,
		Protocol:         *protocol,
		Port:             *port,
	}
	totalTrafficAlert := NewTotalTrafficAlert(*traffic, time.Duration(*duration)*time.Second, ConsoleNotification)
	app := NewApplication(config, HTTPTrafficFilter, totalTrafficAlert)
	app.Run()
}
