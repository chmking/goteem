package state

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/chmking/horde"
	pb "github.com/chmking/horde/protobuf/private"
)

var _ = Describe("StateMachine", func() {
	var sm StateMachine

	Describe("Stopping", func() {
		Context("when the state is UNKNOWN", func() {
			BeforeEach(func() {
				sm.state = pb.Status_UNKNOWN
			})

			It("returns ErrStatusUnknown", func() {
				err := sm.Stopping()
				Expect(err).To(Equal(horde.ErrStatusUnknown))
			})

			It("leaves state UNKNOWN", func() {
				sm.Stopping()
				Expect(sm.State()).To(Equal(pb.Status_UNKNOWN))
			})
		})

		Context("when the state is IDLE", func() {
			BeforeEach(func() {
				sm.state = pb.Status_IDLE
			})

			It("does not return an error", func() {
				err := sm.Stopping()
				Expect(err).To(BeNil())
			})

			It("leaves state IDLE", func() {
				sm.Stopping()
				Expect(sm.State()).To(Equal(pb.Status_IDLE))
			})
		})

		Context("when the state is SCALING", func() {
			BeforeEach(func() {
				sm.state = pb.Status_SCALING
			})

			It("does not return an error", func() {
				err := sm.Stopping()
				Expect(err).To(BeNil())
			})

			It("switches to STOPPING", func() {
				sm.Stopping()
				Expect(sm.State()).To(Equal(pb.Status_STOPPING))
			})
		})

		Context("when the state is RUNNING", func() {
			BeforeEach(func() {
				sm.state = pb.Status_RUNNING
			})

			It("does not return an error", func() {
				err := sm.Stopping()
				Expect(err).To(BeNil())
			})

			It("switches to STOPPING", func() {
				sm.Stopping()
				Expect(sm.State()).To(Equal(pb.Status_STOPPING))
			})
		})

		Context("when the state is STOPPING", func() {
			BeforeEach(func() {
				sm.state = pb.Status_STOPPING
			})

			It("does not return an error", func() {
				err := sm.Stopping()
				Expect(err).To(BeNil())
			})

			It("leaves state STOPPING", func() {
				sm.Stopping()
				Expect(sm.State()).To(Equal(pb.Status_STOPPING))
			})
		})

		Context("when the state is QUITTING", func() {
			BeforeEach(func() {
				sm.state = pb.Status_QUITTING
			})

			It("returns ErrStatusQuitting", func() {
				err := sm.Stopping()
				Expect(err).To(Equal(horde.ErrStatusQuitting))
			})

			It("leaves state QUITTING", func() {
				sm.Stopping()
				Expect(sm.State()).To(Equal(pb.Status_QUITTING))
			})
		})

		Context("when the state is UNEXPECTED", func() {
			BeforeEach(func() {
				sm.state = 42
			})

			It("returns ErrStatusUnexpected", func() {
				err := sm.Stopping()
				Expect(err).To(Equal(horde.ErrStatusUnexpected))
			})

			It("leaves state UNEXPECTED", func() {
				sm.Stopping()
				Expect(sm.State()).To(BeNumerically("==", 42))
			})
		})
	})
})
