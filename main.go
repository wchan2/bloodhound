package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

type Application struct {
	applicationConfig Config
	trafficFilter     TrafficFilter
	alert             Threshold
}

type Config struct {
	Protocol string
	Port     string
}

func NewApplication(config Config, filter TrafficFilter, threshold Threshold) *Application {
	return &Application{
		applicationConfig: config,
		trafficFilter:     filter,
		alert:             threshold,
	}
}

func (a *Application) Run() {
	handle, err := pcap.OpenLive("en0", 1024, false, 1*time.Second)
	if err != nil {
		log.Fatalf("Unable to ", err.Error())
	}
	defer handle.Close()
	handle.SetBPFFilter(fmt.Sprintf("%s port %s"))
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		if !a.trafficFilter.Filter(packet) {
			a.alert.Check(Event{Time: time.Now()})
		}
	}
}

func main() {
	totalTrafficThreshold := NewTotalTrafficThreshold(100, 30*time.Second, ConsoleAlert)
	app := NewApplication(Config{Protocol: "tcp", Port: "80"}, HTTPTrafficFilter, totalTrafficThreshold)
	app.Run()
}
