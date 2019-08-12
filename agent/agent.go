package agent

import (
	"context"
	"net"
	"sync"

	"github.com/chmking/horde"
	"github.com/chmking/horde/session"
	grpc "google.golang.org/grpc"
)

type Session interface {
	Scale(count int32, rate float64, wait int64, cb session.Callback)
	Stop(cb session.Callback)
}

func New(config horde.Config) *Agent {
	return &Agent{
		Session: &session.Session{},
		Status:  horde.Status_IDLE,
	}
}

type Agent struct {
	Session Session
	Status  horde.Status
	server  *grpc.Server
	mtx     sync.Mutex
}

func (a *Agent) SafeStatus() horde.Status {
	a.mtx.Lock()
	defer a.mtx.Unlock()
	return a.Status
}

func (a *Agent) Listen(address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	a.server = grpc.NewServer()
	horde.RegisterAgentServer(a.server, a)
	return a.server.Serve(lis)
}

func (a *Agent) Start(ctx context.Context, req *horde.StartRequest) (*horde.StartResponse, error) {
	switch a.Status {
	case horde.Status_IDLE:
		fallthrough
	case horde.Status_SCALING:
		fallthrough
	case horde.Status_RUNNING:
		a.mtx.Lock()
		a.Status = horde.Status_SCALING
		a.mtx.Unlock()

		a.Session.Scale(req.Users, req.Rate, req.Wait, a.onScaled)
	case horde.Status_STOPPING:
		return nil, horde.ErrStatusStopping
	case horde.Status_QUITTING:
		return nil, horde.ErrStatusQuitting
	}

	return &horde.StartResponse{}, nil
}

func (a *Agent) Stop(ctx context.Context, req *horde.StopRequest) (*horde.StopResponse, error) {
	switch a.Status {
	case horde.Status_IDLE:
		// no-op
	case horde.Status_SCALING:
		fallthrough
	case horde.Status_RUNNING:
		a.mtx.Lock()
		a.Status = horde.Status_STOPPING
		a.mtx.Unlock()

		a.Session.Stop(a.onStopped)
	case horde.Status_STOPPING:
		// no-op
	case horde.Status_QUITTING:
		// no-op
	}

	return &horde.StopResponse{}, nil
}

func (a *Agent) Quit(ctx context.Context, req *horde.QuitRequest) (*horde.QuitResponse, error) {
	defer func() {
		if a.server != nil {
			a.server.Stop()
		}
	}()

	// Regardless of current state, the agent is always switched to
	// QUITTING before exit to deter other requests.
	a.mtx.Lock()
	a.Status = horde.Status_QUITTING
	a.mtx.Unlock()

	return &horde.QuitResponse{}, nil
}

func (a *Agent) onScaled() {
	a.mtx.Lock()
	a.Status = horde.Status_RUNNING
	a.mtx.Unlock()
}

func (a *Agent) onStopped() {
	a.mtx.Lock()
	a.Status = horde.Status_IDLE
	a.mtx.Unlock()
}
