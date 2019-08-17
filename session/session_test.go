package session_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/chmking/horde/session"
)

var _ = Describe("Session", func() {
	var (
		session *Session
		order   ScaleOrder
		done    chan struct{}
	)

	BeforeEach(func() {
		session = &Session{}
		order = ScaleOrder{}
		done = make(chan struct{}, 1)
	})

	Describe("Scale", func() {
		It("calls the Callback when scaled", func() {
			called := false

			session.Scale(order, func() {
				called = true
				close(done)
			})

			<-done

			Expect(called).To(BeTrue())
		})
	})

	Describe("Stop", func() {
		It("calls the Callback when stopped", func() {
			called := false
			session.Stop(func() { called = true })
			Expect(called).To(BeTrue())
		})
	})
})
