package agent

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	"github.com/chmking/horde"
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
		Status:  pb.Status_IDLE,
	}
}

type Agent struct {
	Session Session
	Status  pb.Status
	server  *grpc.Server
	mtx     sync.Mutex
}

func (a *Agent) SafeStatus() pb.Status {
	a.mtx.Lock()
	defer a.mtx.Unlock()
	return a.Status
}

func (a *Agent) Dial(ctx context.Context, address string) {
	conn, err := grpc.Dial("127.0.0.1:5557",
		grpc.WithBackoffMaxDelay(time.Second),
		grpc.WithInsecure(),
		grpc.WithBlock())
	if err != nil {
		log.Fatal(err)
		return
	}

	client := pb.NewManagerClient(conn)

	req := &pb.RegisterRequest{
		Host: "127.0.0.1",
		Port: "5558",
	}

	_, err = client.Register(ctx, req)
	if err != nil {
		return
	}
}

func (a *Agent) Listen(address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	a.server = grpc.NewServer()
	pb.RegisterAgentServer(a.server, a)
	log.Printf("Listening for private connection on %s", address)
	return a.server.Serve(lis)
}

func (a *Agent) Start(ctx context.Context, req *pb.StartRequest) (*pb.StartResponse, error) {
	switch a.Status {
	case pb.Status_IDLE:
		fallthrough
	case pb.Status_SCALING:
		fallthrough
	case pb.Status_RUNNING:
		a.mtx.Lock()
		a.Status = pb.Status_SCALING
		a.mtx.Unlock()

		a.Session.Scale(req.Users, req.Rate, req.Wait, a.onScaled)
	case pb.Status_STOPPING:
		return nil, horde.ErrStatusStopping
	case pb.Status_QUITTING:
		return nil, horde.ErrStatusQuitting
	}

	return &pb.StartResponse{}, nil
}

func (a *Agent) Stop(ctx context.Context, req *pb.StopRequest) (*pb.StopResponse, error) {
	switch a.Status {
	case pb.Status_IDLE:
		// no-op
	case pb.Status_SCALING:
		fallthrough
	case pb.Status_RUNNING:
		a.mtx.Lock()
		a.Status = pb.Status_STOPPING
		a.mtx.Unlock()

		a.Session.Stop(a.onStopped)
	case pb.Status_STOPPING:
		// no-op
	case pb.Status_QUITTING:
		// no-op
	}

	return &pb.StopResponse{}, nil
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
	a.Status = pb.Status_QUITTING
	a.mtx.Unlock()

	return &pb.QuitResponse{}, nil
}

func (a *Agent) Heartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	log.Println("Received a private.Heartbeat request")
	return &pb.HeartbeatResponse{}, nil
}

func (a *Agent) onScaled() {
	a.mtx.Lock()
	a.Status = pb.Status_RUNNING
	a.mtx.Unlock()
}

func (a *Agent) onStopped() {
	a.mtx.Lock()
	a.Status = pb.Status_IDLE
	a.mtx.Unlock()
}
