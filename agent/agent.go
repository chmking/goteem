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

func (a *Agent) Scale(ctx context.Context, req *pb.ScaleRequest) (*pb.ScaleResponse, error) {
	log.Println("Received private.ScaleRequest")
	return &pb.ScaleResponse{}, nil
}

func (a *Agent) Stop(ctx context.Context, req *pb.StopRequest) (*pb.StopResponse, error) {
	log.Println("Received private.StopRequest")
	return &pb.StopResponse{}, nil
}

func (a *Agent) Quit(ctx context.Context, req *pb.QuitRequest) (*pb.QuitResponse, error) {
	log.Println("Received private.QuitRequest")
	return &pb.QuitResponse{}, nil
}

func (a *Agent) Heartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
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
