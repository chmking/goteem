package state

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/chmking/horde"
	pb "github.com/chmking/horde/protobuf/public"
)

var _ = Describe("StateMachine", func() {
	var sm StateMachine

	Describe("Idle", func() {
		Context("when the state is IDLE", func() {
			BeforeEach(func() {
				sm.state = pb.Status_STATUS_IDLE
			})

			It("does not return an error", func() {
				err := sm.Idle()
				Expect(err).To(BeNil())
			})

			It("leaves state IDLE", func() {
				sm.Idle()
				Expect(sm.State()).To(Equal(pb.Status_STATUS_IDLE))
			})
		})

		Context("when the state is UNKNOWN", func() {
			BeforeEach(func() {
				sm.state = pb.Status_STATUS_UNKNOWN
			})

			It("does not return an error", func() {
				err := sm.Idle()
				Expect(err).To(BeNil())
			})

			It("switches to IDLE", func() {
				sm.Idle()
				Expect(sm.State()).To(Equal(pb.Status_STATUS_IDLE))
			})
		})

		Context("when the state is SCALING", func() {
			BeforeEach(func() {
				sm.state = pb.Status_STATUS_SCALING
			})

			It("returns ErrStatusScaling", func() {
				err := sm.Idle()
				Expect(err).To(Equal(horde.ErrStatusScaling))
			})

			It("leaves state SCALING", func() {
				sm.Idle()
				Expect(sm.State()).To(Equal(pb.Status_STATUS_SCALING))
			})
		})

		Context("when the state is RUNNING", func() {
			BeforeEach(func() {
				sm.state = pb.Status_STATUS_RUNNING
			})

			It("returns ErrStatusRunning", func() {
				err := sm.Idle()
				Expect(err).To(Equal(horde.ErrStatusRunning))
			})

			It("leaves state RUNNING", func() {
				sm.Idle()
				Expect(sm.State()).To(Equal(pb.Status_STATUS_RUNNING))
			})
		})

		Context("when the state is STOPPING", func() {
			BeforeEach(func() {
				sm.state = pb.Status_STATUS_STOPPING
			})

			It("does not return an error", func() {
				err := sm.Idle()
				Expect(err).To(BeNil())
			})

			It("switches to IDLE", func() {
				sm.Idle()
				Expect(sm.State()).To(Equal(pb.Status_STATUS_IDLE))
			})
		})

		Context("when the state is QUITTING", func() {
			BeforeEach(func() {
				sm.state = pb.Status_STATUS_QUITTING
			})

			It("returns ErrStatusQuitting", func() {
				err := sm.Idle()
				Expect(err).To(Equal(horde.ErrStatusQuitting))
			})

			It("leaves state QUITTING", func() {
				sm.Idle()
				Expect(sm.State()).To(Equal(pb.Status_STATUS_QUITTING))
			})
		})

		Context("when the state is UNEXPECTED", func() {
			BeforeEach(func() {
				sm.state = 42
			})

			It("returns ErrStatusUnexpected", func() {
				err := sm.Idle()
				Expect(err).To(Equal(horde.ErrStatusUnexpected))
			})

			It("leaves state UNEXPECTED", func() {
				sm.Idle()
				Expect(sm.State()).To(BeNumerically("==", 42))
			})
		})
	})
})
