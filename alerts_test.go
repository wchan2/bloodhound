package main_test

import (
	"bytes"
	"log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/wchan2/bloodhound"
)

var _ = Describe(`ConsoleAlert`, func() {
	Describe(`#Send`, func() {
		var buf *bytes.Buffer

		BeforeEach(func() {
			buf = new(bytes.Buffer)
			log.SetOutput(buf)
			log.SetFlags(0)
		})

		JustBeforeEach(func() {
			ConsoleAlert.Send("alert message")
		})

		It(`sends an alert to the console`, func() {
			Expect(buf.String()).To(Equal("alert message\n"))
		})
	})
})
