package main

import "log"

type Alert interface {
	Send(string)
}

var ConsoleAlert = NewAlert(func(message string) {
	log.Printf(message)
})

type alert struct {
	sendFn SendFn
}

type SendFn func(string)

func NewAlert(fn SendFn) *alert {
	return &alert{sendFn: fn}
}

func (a *alert) Send(message string) {
	a.sendFn(message)
}
