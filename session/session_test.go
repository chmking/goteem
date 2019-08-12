package session_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/chmking/horde/session"
)

var _ = Describe("Session", func() {
	var (
		session *Session
		done    chan struct{}
	)

	BeforeEach(func() {
		session = &Session{}
		done = make(chan struct{}, 1)
	})

	Describe("Scale", func() {
		It("calls the Callback when scaled", func() {
			called := false

			session.Scale(0, 0, 0, func() {
				if done != nil {
					called = true
					close(done)
				}
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
