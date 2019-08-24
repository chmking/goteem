package eventloop_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/chmking/horde/eventloop"
)

var _ = Describe("Eventloop", func() {
	var (
		el *EventLoop
	)

	BeforeEach(func() {
		el = New()
	})

	Describe("New", func() {
		It("returns a new EventLoop", func() {
			value := New()
			Expect(value).NotTo(BeNil())
		})
	})

	Describe("Append", func() {
		It("adds the Event to the EventLoop", func() {
			example := func(param string) (out string) {
				el.Append(func() {
					out = param
				})

				return
			}

			value := example("foo")
			Expect(value).To(Equal("foo"))
		})
	})
})
