package main

import (
	"flag"
	"time"
)

var (
	iface    = flag.String("interface", "en0", "Network interface to listen for packets")
	protocol = flag.String("protocol", "tcp", "Protocol to listen for packets")
	port     = flag.String("port", "80", "Port to listen to for packets")

	monitor  = flag.Int("monitor", 10, "Monitoring duration in seconds to which to send a summary")
	duration = flag.Int("duration", 30, "Duration in seconds that")
	traffic  = flag.Int("traffic", 100, "Traffic amount that should trigger an alert")
)

func main() {
	config := Config{
		NetworkInterface: *iface,
		Protocol:         *protocol,
		Port:             *port,
	}
	totalTrafficAlert := NewTotalTrafficAlert(*traffic, time.Duration(*duration)*time.Second, ConsoleNotification)
	trafficMonitor := NewSummaryStatsTrafficMonitor(time.Duration(*monitor)*time.Second, ConsoleNotification)
	app := NewApplication(config, HTTPTrafficFilter, trafficMonitor, totalTrafficAlert)
	app.Run()
}
