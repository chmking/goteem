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
				mockClient.EXPECT().Scale(gomock.Any(), gomock.Any()).Return(&private.ScaleResponse{}, nil).AnyTimes()
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

			Context("and the manager is IDLE", func() {
				BeforeEach(func() {
					mockStateMachine.EXPECT().State().Return(public.Status_STATUS_IDLE).AnyTimes()
					mockStateMachine.EXPECT().Scaling().Return(nil).AnyTimes()
				})

				It("does not return an error", func() {
					err := manager.Start(1, 1)
					Expect(err).To(BeNil())
				})
			})

			Context("and the manager is SCALING", func() {
				BeforeEach(func() {
					mockStateMachine.EXPECT().State().Return(public.Status_STATUS_SCALING).AnyTimes()
					mockStateMachine.EXPECT().Scaling().Return(nil).AnyTimes()
				})

				It("does not return an error", func() {
					err := manager.Start(1, 1)
					Expect(err).To(BeNil())
				})
			})

			Context("and the manager is RUNNING", func() {
				BeforeEach(func() {
					mockStateMachine.EXPECT().State().Return(public.Status_STATUS_RUNNING).AnyTimes()
					mockStateMachine.EXPECT().Scaling().Return(nil).AnyTimes()
				})

				It("does not return an error", func() {
					err := manager.Start(1, 1)
					Expect(err).To(BeNil())
				})
			})
		})
	})
})
