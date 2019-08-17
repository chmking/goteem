package state

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/chmking/horde"
	pb "github.com/chmking/horde/protobuf/public"
)

var _ = Describe("StateMachine", func() {
	var sm StateMachine

	Describe("Scaling", func() {
		Context("when the state is UNKNOWN", func() {
			BeforeEach(func() {
				sm.state = pb.Status_STATUS_UNKNOWN
			})

			It("returns ErrStatusUnknown", func() {
				err := sm.Scaling()
				Expect(err).To(Equal(horde.ErrStatusUnknown))
			})

			It("leaves state UNKNOWN", func() {
				sm.Scaling()
				Expect(sm.State()).To(Equal(pb.Status_STATUS_UNKNOWN))
			})
		})

		Context("when the state is IDLE", func() {
			BeforeEach(func() {
				sm.state = pb.Status_STATUS_IDLE
			})

			It("does not return an error", func() {
				err := sm.Scaling()
				Expect(err).To(BeNil())
			})

			It("switches to SCALING", func() {
				sm.Scaling()
				Expect(sm.State()).To(Equal(pb.Status_STATUS_SCALING))
			})
		})

		Context("when the state is SCALING", func() {
			BeforeEach(func() {
				sm.state = pb.Status_STATUS_SCALING
			})

			It("does not return an error", func() {
				err := sm.Scaling()
				Expect(err).To(BeNil())
			})

			It("leaves state SCALING", func() {
				sm.Scaling()
				Expect(sm.State()).To(Equal(pb.Status_STATUS_SCALING))
			})
		})

		Context("when the state is RUNNING", func() {
			BeforeEach(func() {
				sm.state = pb.Status_STATUS_RUNNING
			})

			It("does not return an error", func() {
				err := sm.Scaling()
				Expect(err).To(BeNil())
			})

			It("switches to SCALING", func() {
				sm.Scaling()
				Expect(sm.State()).To(Equal(pb.Status_STATUS_SCALING))
			})
		})

		Context("when the state is STOPPING", func() {
			BeforeEach(func() {
				sm.state = pb.Status_STATUS_STOPPING
			})

			It("returns ErrStatusStopping", func() {
				err := sm.Scaling()
				Expect(err).To(Equal(horde.ErrStatusStopping))
			})

			It("leaves state STOPPING", func() {
				sm.Scaling()
				Expect(sm.State()).To(Equal(pb.Status_STATUS_STOPPING))
			})
		})

		Context("when the state is QUITTING", func() {
			BeforeEach(func() {
				sm.state = pb.Status_STATUS_QUITTING
			})

			It("returns ErrStatusQuitting", func() {
				err := sm.Scaling()
				Expect(err).To(Equal(horde.ErrStatusQuitting))
			})

			It("leaves state QUITTING", func() {
				sm.Scaling()
				Expect(sm.State()).To(Equal(pb.Status_STATUS_QUITTING))
			})
		})

		Context("when the state is UNEXPECTED", func() {
			BeforeEach(func() {
				sm.state = 42
			})

			It("returns ErrStatusUnexpected", func() {
				err := sm.Scaling()
				Expect(err).To(Equal(horde.ErrStatusUnexpected))
			})

			It("leaves state UNEXPECTED", func() {
				sm.Scaling()
				Expect(sm.State()).To(BeNumerically("==", 42))
			})
		})
	})
})
