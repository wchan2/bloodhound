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

### Running the application

Running the application with the below command will require building it in [this section](#building).

Note: the `sudo` may be required to allow the application to listen to the specified network interface.

```bash
sudo ./bloodhound
```

### Flags

Some flags that can be used to customize the application at runtime.

```bash
# network flags
- interface - Network interface to listen for packets
	- default: "en0"
- protocol - Protocol to listen for packets
	- default: "tcp"
- port - Port to listen to for packets
	- default: "80"

# monitoring and alerting flags
- monitor - Monitoring duration in seconds to which to send a summary
	- default: 10
- duration - Duration in seconds that
	- default: 30
- traffic - Traffic amount that should trigger an alert
	- default: 100
```

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

## Design

Below are some of the extensible components, namely interfaces and what their responsibilities are. Under each component are a list of pre-existing components that implements the respective interface.

### Interfaces and Implementations

Components that can be extended or customized to be used in the application.

- `TrafficFilter` decides what messages to filter out and keep
	- `HTTPTrafficFilter` filters all traffic that are not HTTP traffic
- `Monitor` monitors traffic
	- `TrafficMonitor` generates statistical summaries for traffic received and sent
- `Alert` evaluates whether an event surpasses the threshold or reverts to normal
	- `TotalTrafficAlert` keeps track of the total number of events in a given time window
- `Notification` that determines when to alert
	- `ConsoleNotification` alerts to the console

## Domain Messages

Messages that are passed from one component to another.

- `Event` represents a network event with fields such as status, payload, sender, destination, etc
- `TrafficStatistics` has fields for different traffic statistics such as average payload size and total payload size

## Application

Application that listens to network traffic and passes it through a filter, a monitor, a threshold, and eventually an alert if traffic surpasses the threshold.

- `Application` is composed of the different interfaces, namely the `TrafficFilter`, `Monitor`, `Alert`, and `Notification` to allow custom components to filter for relevant traffic, monitor the filtered traffic, and alert when when the traffic surpasses some threshold

## License

bloodhound is released under the [MIT License](https://opensource.org/licenses/MIT).
