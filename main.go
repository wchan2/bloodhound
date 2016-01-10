package main

import (
	"flag"
	"log"
	"time"
)

var (
	iface    string
	protocol string
	port     string

	monitor  int
	duration int
	traffic  int
)

func init() {
	flag.StringVar(&iface, "interface", "en0", "Network interface to listen for packets")
	flag.StringVar(&protocol, "protocol", "tcp", "Protocol to listen for packets")
	flag.StringVar(&port, "port", "80", "Port to listen to for packets")
	flag.IntVar(&monitor, "monitor", 10, "Monitoring duration in seconds to which to send a summary")
	flag.IntVar(&duration, "duration", 30, "Duration in seconds that")
	flag.IntVar(&traffic, "traffic", 100, "Traffic amount that should trigger an alert")
}

func main() {
	flag.Parse()
	config := Config{
		NetworkInterface: iface,
		Protocol:         protocol,
		Port:             port,
	}
	totalTrafficAlert := NewTotalTrafficAlert(traffic, time.Duration(duration)*time.Second, ConsoleNotification)
	trafficMonitor := NewSummaryStatsTrafficMonitor(time.Duration(monitor)*time.Second, ConsoleNotification)

	app := NewApplication(config, HTTPTrafficFilter, trafficMonitor, totalTrafficAlert)

	log.Printf("Listening to %s traffic on the %s interface on port %s", protocol, iface, port)
	log.Printf("Monitoring traffic; will alert if traffic surpasses %d per second", traffic)
	app.Run()
}
