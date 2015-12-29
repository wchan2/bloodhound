# bloodhound

An extensible packet sniffer application that can filter and monitor network traffic. It also has capabilities of producing alerts.

## Dependencies

### Operating System Dependencies

Depending on the operating system, the install instructions will be different. 

- [libpcap](http://www.tcpdump.org/#latest-release) for monitoring network traffic on the operating system

### Go Dependencies

The below dependencies are managed by the [Godeps](http://github.com/tools/godep) and will require godeps to be installed. Please see the [Godep installation](https://github.com/tools/godep#install) for more instructions.

- [gopacket](https://github.com/google/gopacket) for packet sniffing
- [ginkgo](https://github.com/onsi/ginkgo) for BDD style tests
- [gomega](github.com/onsi/gomega) for matchers used to create assertions in gingko

## Documentation

TBD

## Design

Below are some of the extensible components, namely interfaces and what their responsibilities are. Under each component are a list of pre-existing components that implements the respective interface.

- `TrafficFilter` decides what messages to filter out and keep
	- `HTTPTrafficFilter` filters all traffic that are not HTTP traffic
- `Monitor` monitors traffic
	- `TrafficMonitor` generates statistical summaries for traffic received and sent
- `Alert` that determines when to alert
	- `ConsoleAlert` alerts to the console
- `Threshold` evaluates whether an event surpasses the threshold or reverts to normal
	- `TotalTrafficThreshold` keeps track of the total number of events in a given time window

## Building 

Run the below command in the

```bash
godep go build
```

## Running tests

Run the below command in the directory of the top most directory of the project.

```bash
godep go test ./...
```

## License

bloodhound is released under the [MIT License](https://opensource.org/licenses/MIT).
