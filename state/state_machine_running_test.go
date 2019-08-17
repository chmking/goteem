package state

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/chmking/horde"
	pb "github.com/chmking/horde/protobuf/private"
)

var _ = Describe("StateMachine", func() {
	var sm StateMachine

	Describe("Scaling", func() {
		Context("when the state is UNKNOWN", func() {
			BeforeEach(func() {
				sm.state = pb.Status_UNKNOWN
			})

			It("returns ErrStatusUnknown", func() {
				err := sm.Running()
				Expect(err).To(Equal(horde.ErrStatusUnknown))
			})

			It("leaves state UNKNOWN", func() {
				sm.Running()
				Expect(sm.State()).To(Equal(pb.Status_UNKNOWN))
			})
		})

		Context("when the state is IDLE", func() {
			BeforeEach(func() {
				sm.state = pb.Status_IDLE
			})

			It("returns ErrStatusIdle", func() {
				err := sm.Running()
				Expect(err).To(Equal(horde.ErrStatusIdle))
			})

			It("leaves state UNKNOWN", func() {
				sm.Running()
				Expect(sm.State()).To(Equal(pb.Status_IDLE))
			})
		})

		Context("when the state is SCALING", func() {
			BeforeEach(func() {
				sm.state = pb.Status_SCALING
			})

			It("does not return an error", func() {
				err := sm.Running()
				Expect(err).To(BeNil())
			})

			It("switches to RUNNING", func() {
				sm.Running()
				Expect(sm.State()).To(Equal(pb.Status_RUNNING))
			})
		})

		Context("when the state is RUNNING", func() {
			BeforeEach(func() {
				sm.state = pb.Status_RUNNING
			})

			It("does not return an error", func() {
				err := sm.Running()
				Expect(err).To(BeNil())
			})

			It("leaves state RUNNING", func() {
				sm.Running()
				Expect(sm.State()).To(Equal(pb.Status_RUNNING))
			})
		})

		Context("when the state is STOPPING", func() {
			BeforeEach(func() {
				sm.state = pb.Status_STOPPING
			})

			It("returns ErrStatusStopping", func() {
				err := sm.Running()
				Expect(err).To(Equal(horde.ErrStatusStopping))
			})

			It("leaves state STOPPING", func() {
				sm.Running()
				Expect(sm.State()).To(Equal(pb.Status_STOPPING))
			})
		})

		Context("when the state is QUITTING", func() {
			BeforeEach(func() {
				sm.state = pb.Status_QUITTING
			})

			It("returns ErrStatusQuitting", func() {
				err := sm.Running()
				Expect(err).To(Equal(horde.ErrStatusQuitting))
			})

			It("leaves state QUITTING", func() {
				sm.Running()
				Expect(sm.State()).To(Equal(pb.Status_QUITTING))
			})
		})
	})
})
