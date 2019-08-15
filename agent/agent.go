package agent

import (
	"context"
	"net"
	"sync"

	"github.com/chmking/horde"
	"github.com/chmking/horde/protobuf/private"
	pb "github.com/chmking/horde/protobuf/private"
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
		Status:  private.Status_IDLE,
	}
}

type Agent struct {
	Session Session
	Status  private.Status
	server  *grpc.Server
	mtx     sync.Mutex
}

func (a *Agent) SafeStatus() private.Status {
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
	return a.server.Serve(lis)
}

func (a *Agent) Start(ctx context.Context, req *pb.StartRequest) (*pb.StartResponse, error) {
	switch a.Status {
	case private.Status_IDLE:
		fallthrough
	case private.Status_SCALING:
		fallthrough
	case private.Status_RUNNING:
		a.mtx.Lock()
		a.Status = private.Status_SCALING
		a.mtx.Unlock()

		a.Session.Scale(req.Users, req.Rate, req.Wait, a.onScaled)
	case private.Status_STOPPING:
		return nil, horde.ErrStatusStopping
	case private.Status_QUITTING:
		return nil, horde.ErrStatusQuitting
	}

	return &private.StartResponse{}, nil
}

func (a *Agent) Stop(ctx context.Context, req *pb.StopRequest) (*pb.StopResponse, error) {
	switch a.Status {
	case private.Status_IDLE:
		// no-op
	case private.Status_SCALING:
		fallthrough
	case private.Status_RUNNING:
		a.mtx.Lock()
		a.Status = private.Status_STOPPING
		a.mtx.Unlock()

		a.Session.Stop(a.onStopped)
	case private.Status_STOPPING:
		// no-op
	case private.Status_QUITTING:
		// no-op
	}

	return &private.StopResponse{}, nil
}

func (a *Agent) Quit(ctx context.Context, req *pb.QuitRequest) (*pb.QuitResponse, error) {
	defer func() {
		if a.server != nil {
			a.server.Stop()
		}
	}()

	// Regardless of current state, the agent is always switched to
	// QUITTING before exit to deter other requests.
	a.mtx.Lock()
	a.Status = private.Status_QUITTING
	a.mtx.Unlock()

	return &private.QuitResponse{}, nil
}

func (a *Agent) onScaled() {
	a.mtx.Lock()
	a.Status = private.Status_RUNNING
	a.mtx.Unlock()
}

func (a *Agent) onStopped() {
	a.mtx.Lock()
	a.Status = private.Status_IDLE
	a.mtx.Unlock()
}
