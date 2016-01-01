package main

import "log"

type Notification interface {
	Send(string)
}

var ConsoleNotification = NewNotificationSender(func(message string) {
	log.Printf(message)
})

type NotificationSender struct {
	sendFn SendFn
}

type SendFn func(string)

func NewNotificationSender(fn SendFn) *NotificationSender {
	return &NotificationSender{sendFn: fn}
}

func (n *NotificationSender) Send(message string) {
	n.sendFn(message)
}
