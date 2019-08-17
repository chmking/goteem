package state

import (
	"sync"

	"github.com/chmking/horde"
	pb "github.com/chmking/horde/protobuf/private"
)

type StateMachine struct {
	state pb.Status
	mtx   sync.Mutex
}

func (sm *StateMachine) setState(state pb.Status) {
	sm.mtx.Lock()
	sm.state = state
	sm.mtx.Unlock()
}

func (sm *StateMachine) State() pb.Status {
	sm.mtx.Lock()
	defer sm.mtx.Unlock()
	return sm.state
}

func (sm *StateMachine) Idle() error {
	switch sm.state {
	case pb.Status_IDLE:
		// no-op
		return nil
	case pb.Status_UNKNOWN:
		fallthrough
	case pb.Status_STOPPING:
		sm.setState(pb.Status_IDLE)
		return nil
	case pb.Status_SCALING:
		return horde.ErrStatusScaling
	case pb.Status_RUNNING:
		return horde.ErrStatusRunning
	case pb.Status_QUITTING:
		return horde.ErrStatusQuitting
	default:
		return horde.ErrStatusUnexpected
	}
}

func (sm *StateMachine) Scaling() error {
	switch sm.state {
	case pb.Status_SCALING:
		// no-op
		return nil
	case pb.Status_IDLE:
		fallthrough
	case pb.Status_RUNNING:
		sm.setState(pb.Status_SCALING)
		return nil
	case pb.Status_UNKNOWN:
		return horde.ErrStatusUnknown
	case pb.Status_STOPPING:
		return horde.ErrStatusStopping
	case pb.Status_QUITTING:
		return horde.ErrStatusQuitting
	default:
		return horde.ErrStatusUnexpected
	}
}

func (sm *StateMachine) Running() error {
	switch sm.state {
	case pb.Status_RUNNING:
		// no-op
		return nil
	case pb.Status_SCALING:
		sm.setState(pb.Status_RUNNING)
		return nil
	case pb.Status_IDLE:
		return horde.ErrStatusIdle
	case pb.Status_UNKNOWN:
		return horde.ErrStatusUnknown
	case pb.Status_STOPPING:
		return horde.ErrStatusStopping
	case pb.Status_QUITTING:
		return horde.ErrStatusQuitting
	default:
		return horde.ErrStatusUnexpected
	}
}

func (sm *StateMachine) Stopping() error {
	switch sm.state {
	case pb.Status_IDLE:
		fallthrough
	case pb.Status_STOPPING:
		// no-op
		return nil
	case pb.Status_SCALING:
		fallthrough
	case pb.Status_RUNNING:
		sm.setState(pb.Status_STOPPING)
		return nil
	case pb.Status_UNKNOWN:
		return horde.ErrStatusUnknown
	case pb.Status_QUITTING:
		return horde.ErrStatusQuitting
	default:
		return horde.ErrStatusUnexpected
	}
}

func (sm *StateMachine) Quitting() error {
	switch sm.state {
	case pb.Status_QUITTING:
		// no-op
		return nil
	case pb.Status_IDLE:
		fallthrough
	case pb.Status_STOPPING:
		fallthrough
	case pb.Status_SCALING:
		fallthrough
	case pb.Status_RUNNING:
		sm.setState(pb.Status_QUITTING)
		return nil
	case pb.Status_UNKNOWN:
		return horde.ErrStatusUnknown
	default:
		return horde.ErrStatusUnexpected
	}
}
