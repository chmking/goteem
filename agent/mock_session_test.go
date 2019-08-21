package agent_test

import (
	"github.com/chmking/horde/agent/session"
)

type MockSession struct {
}

func (m *MockSession) Scale(count int32, rate float64, wait int64, cb session.Callback) {
	// no-op
}

func (m *MockSession) Stop(session.Callback) {
	// no-op
}
