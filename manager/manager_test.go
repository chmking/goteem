package manager

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/chmking/horde/protobuf/public"
)

type MockStateMachine struct {
	Current public.Status

	stateMachine
}

func (m *MockStateMachine) State() public.Status {
	return m.Current
}

var _ = Describe("Manager", func() {
	var (
		m   *Manager
		msm *MockStateMachine
	)

	BeforeEach(func() {
		msm = &MockStateMachine{}

		m = &Manager{
			sm: msm,
		}
	})

	Describe("State", func() {
		BeforeEach(func() {
			msm.Current = public.Status_STATUS_RUNNING
		})

		It("returns the manager state", func() {
			value := m.State()
			Expect(value).To(Equal(public.Status_STATUS_RUNNING))
		})
	})
})
