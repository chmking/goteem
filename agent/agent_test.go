package agent

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/chmking/horde"
	gomock "github.com/golang/mock/gomock"
)

var _ = Describe("Agent", func() {
	var (
		agent  *Agent
		config horde.Config

		mockCtrl         *gomock.Controller
		mockSession      *MockSession
		mockStateMachine *MockStateMachine
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockSession = NewMockSession(mockCtrl)
		mockStateMachine = NewMockStateMachine(mockCtrl)

		agent = New(config)
		agent.Session = mockSession
		agent.StateMachine = mockStateMachine
	})

	Describe("New", func() {
		It("constructs a new agent", func() {
			value := New(horde.Config{})
			Expect(value).ToNot(BeNil())
		})
	})

	Describe("Scale", func() {
		var orders Orders

		Context("when the agent is in an invalid state", func() {
			BeforeEach(func() {
				mockStateMachine.EXPECT().Scaling().Return(errors.New("foo")).AnyTimes()
			})

			It("returns an error", func() {
				err := agent.Scale(orders)
				Expect(err).To(Equal(errors.New("foo")))
			})
		})

		Context("when the agent is in a valid state", func() {
			BeforeEach(func() {
				mockStateMachine.EXPECT().Scaling().Return(nil).AnyTimes()
				mockSession.EXPECT().Scale(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
			})

			It("does not return an error", func() {
				err := agent.Scale(orders)
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("Stop", func() {
		Context("when the agent is in an invalid state", func() {
			BeforeEach(func() {
				mockStateMachine.EXPECT().Stopping().Return(errors.New("foo")).AnyTimes()
			})

			It("returns an error", func() {
				err := agent.Stop()
				Expect(err).To(Equal(errors.New("foo")))
			})
		})

		Context("when the agent is in a valid state", func() {
			BeforeEach(func() {
				mockStateMachine.EXPECT().Stopping().Return(nil).AnyTimes()
				mockSession.EXPECT().Stop(gomock.Any()).AnyTimes()
			})

			It("does not return an error", func() {
				err := agent.Stop()
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("onScaled", func() {
		var running bool

		BeforeEach(func() {
			running = false
			mockStateMachine.EXPECT().Running().Do(func() { running = true }).AnyTimes()
		})

		It("sets the agent to RUNNING", func() {
			agent.onScaled()
			Expect(running).To(BeTrue())
		})
	})

	Describe("onStopped", func() {
		var idle bool

		BeforeEach(func() {
			idle = false
			mockStateMachine.EXPECT().Idle().Do(func() { idle = true }).AnyTimes()
		})

		It("sets the agent to RUNNING", func() {
			agent.onStopped()
			Expect(idle).To(BeTrue())
		})
	})
})
