package agent

import (
	"context"
	"net"

	"github.com/chmking/horde"
	"github.com/chmking/horde/session"
	grpc "google.golang.org/grpc"
)

type Session interface {
	Scale(count int32, rate float64, wait int64)
	Stop()
}

func New(config horde.Config) *Agent {
	return &Agent{
		status:  horde.Status_IDLE,
		session: &session.Session{},
	}
}

type Agent struct {
	session Session
	status  horde.Status
	server  *grpc.Server
}

func (a *Agent) Status() horde.Status {
	return a.status
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
	switch a.status {
	case horde.Status_IDLE:
		fallthrough
	case horde.Status_RUNNING:
		a.status = horde.Status_RUNNING
		a.session.Scale(req.Users, req.Rate, req.Wait)
	case horde.Status_STOPPING:
		return nil, horde.ErrStatusStopping
	case horde.Status_QUITTING:
		return nil, horde.ErrStatusQuitting
	}

	return &horde.StartResponse{}, nil
}

func (a *Agent) Stop(ctx context.Context, req *horde.StopRequest) (*horde.StopResponse, error) {
	switch a.status {
	case horde.Status_IDLE:
		// no-op
	case horde.Status_RUNNING:
		a.status = horde.Status_STOPPING
		a.session.Stop()
		a.status = horde.Status_IDLE
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
	a.status = horde.Status_QUITTING

	return &horde.QuitResponse{}, nil
}
