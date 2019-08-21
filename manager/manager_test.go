package manager_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/chmking/horde/manager"
	"github.com/chmking/horde/manager/registry"
	"github.com/chmking/horde/protobuf/private"
	"github.com/chmking/horde/protobuf/public"
	gomock "github.com/golang/mock/gomock"
)

var _ = Describe("Manager", func() {
	var (
		manager          *Manager
		mockCtrl         *gomock.Controller
		mockClient       *MockAgentClient
		mockRegistry     *MockRegistry
		mockStateMachine *MockStateMachine
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockClient = NewMockAgentClient(mockCtrl)
		mockRegistry = NewMockRegistry(mockCtrl)
		mockStateMachine = NewMockStateMachine(mockCtrl)

		manager = &Manager{
			Registry:     mockRegistry,
			StateMachine: mockStateMachine,
		}
	})

	Describe("Start", func() {
		Context("when there are no registered agents", func() {
			BeforeEach(func() {
				mockRegistry.EXPECT().GetActive().Return(nil).AnyTimes()
			})

			It("returns ErrNoActiveAgents", func() {
				err := manager.Start(1, 1)
				Expect(err).To(Equal(ErrNoActiveAgents))
			})
		})

		Context("when there is at least one registered agent", func() {
			BeforeEach(func() {
				mockRegistry.EXPECT().GetActive().Return([]registry.Registration{
					registry.Registration{Client: mockClient},
				}).AnyTimes()
			})

			Context("and the manager is in an invalid state", func() {
				BeforeEach(func() {
					mockStateMachine.EXPECT().State().Return(public.Status_STATUS_QUITTING).AnyTimes()
					mockStateMachine.EXPECT().Scaling().Return(errors.New("foo")).AnyTimes()
				})

				It("returns an error", func() {
					err := manager.Start(1, 1)
					Expect(err).To(Equal(errors.New("foo")))
				})
			})

			Context("and the manager is in a valid state", func() {
				BeforeEach(func() {
					mockStateMachine.EXPECT().State().Return(public.Status_STATUS_IDLE).AnyTimes()
					mockStateMachine.EXPECT().Scaling().Return(nil).AnyTimes()
				})

				Context("and the agent scale does not return an error", func() {
					BeforeEach(func() {
						mockClient.EXPECT().Scale(gomock.Any(), gomock.Any()).Return(&private.ScaleResponse{}, nil)
					})

					It("does not return an error", func() {
						err := manager.Start(1, 1)
						Expect(err).To(BeNil())
					})
				})

				Context("and the agent scale returns an error", func() {
					BeforeEach(func() {
						mockClient.EXPECT().Scale(gomock.Any(), gomock.Any()).Return(&private.ScaleResponse{}, errors.New("foo"))
					})

					It("does not return an error", func() {
						err := manager.Start(1, 1)
						Expect(err).To(BeNil())
					})
				})
			})
		})
	})

	Describe("Stop", func() {
		BeforeEach(func() {
			mockRegistry.EXPECT().GetAll().Return(nil).AnyTimes()
		})

		Context("when the manager is in an invalid state", func() {
			BeforeEach(func() {
				mockStateMachine.EXPECT().Stopping().Return(errors.New("foo")).AnyTimes()
			})

			It("returns an error", func() {
				err := manager.Stop()
				Expect(err).To(Equal(errors.New("foo")))
			})
		})

		Context("when the manager is in a valid state", func() {
			BeforeEach(func() {
				mockStateMachine.EXPECT().Stopping().Return(nil).AnyTimes()
			})

			Context("and there is no registered agent", func() {
				BeforeEach(func() {
					mockRegistry.EXPECT().GetAll().Return(nil).AnyTimes()
				})

				It("does not return an error", func() {
					err := manager.Stop()
					Expect(err).To(BeNil())
				})
			})

			Context("and there is at least one registered agent", func() {
				BeforeEach(func() {
					mockRegistry.EXPECT().GetAll().Return([]registry.Registration{
						registry.Registration{Client: mockClient},
					}).AnyTimes()
				})

				It("calls Stop on the agents", func() {
					mockClient.EXPECT().Stop(gomock.Any(), gomock.Any()).Return(&private.StopResponse{}, nil)
				})

				It("does not return an error", func() {
					err := manager.Stop()
					Expect(err).To(BeNil())
				})

				Context("and the agent reports an error", func() {
					BeforeEach(func() {
						mockClient.EXPECT().Stop(gomock.Any(), gomock.Any()).Return(&private.StopResponse{}, errors.New("foo"))
					})

					It("does not return an error", func() {
						err := manager.Stop()
						Expect(err).To(BeNil())
					})
				})
			})
		})
	})
})
