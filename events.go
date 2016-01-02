package main

import "time"

type Event struct {
	Sender      string
	Identifier  string
	User        string
	Time        time.Time
	Destination string
	Status      string
	Payload     []byte
}
