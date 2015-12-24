package main

func main() {
	// components
	// - traffic listener - decides what kind of packets to listen to
	// - a component that reads the thresholds and calls the proper alert
	// - alerting component - deciding when to alert and when to recover
	// - write a summary layer - statistics of most popular site
	// - write a presentational layer that outputs the crossing of the threshold

	// sniffer := NewPacketSniffer()
	// sniffer.SetTrafficListener(NewHTTPTrafficListener()) // filters network traffic
	// sniffer.SetTrafficMonitor(NewTrafficMonitor())
	// // sniffer.SetAlert(NewConsoleAlert.Watch(NewTrafficMonitor()))
	// err := sniffer.Run()
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
}
