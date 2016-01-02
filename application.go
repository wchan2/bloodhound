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
	trafficMonitor    TrafficMonitor
	alert             Alert
}

type Config struct {
	NetworkInterface string
	Protocol         string
	Port             string
}

func NewApplication(config Config, filter TrafficFilter, monitor TrafficMonitor, alert Alert) *Application {
	return &Application{
		applicationConfig: config,
		trafficFilter:     filter,
		trafficMonitor:    monitor,
		alert:             alert,
	}
}

func (a *Application) Run() {
	handle, err := pcap.OpenLive(a.applicationConfig.NetworkInterface, 1024, false, 1*time.Second)
	if err != nil {
		log.Fatalf("Unable to ", err.Error())
	}
	defer handle.Close()
	defer a.trafficMonitor.Stop()
	handle.SetBPFFilter(fmt.Sprintf("%s port %s", a.applicationConfig.Protocol, a.applicationConfig.Port))
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		if event, filtered := a.trafficFilter.Filter(packet); !filtered {
			a.trafficMonitor.Monitor(event)
			a.alert.Check(event)
		}
	}
}
