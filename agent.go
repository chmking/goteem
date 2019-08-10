package goteem

import (
	"context"
	"net"

	grpc "google.golang.org/grpc"
)

type Task struct {
	Name   string
	Weight int
	Func   func(ctx context.Context)
}

type Behavior struct {
	Tasks []*Task
}

type Agent struct {
	Behavior *Behavior

	server *grpc.Server
}

func (a *Agent) Listen(address string) error {
	// Open socket
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	// Start gRPC server
	a.server = grpc.NewServer()
	RegisterAgentServer(a.server, a)
	return a.server.Serve(lis)
}

func (a *Agent) Teem(ctx context.Context, req *TeemRequest) (*TeemResponse, error) {
	return &TeemResponse{}, nil
}

func (a *Agent) Quit(ctx context.Context, req *QuitRequest) (*QuitResponse, error) {
	defer func() {
		a.server.Stop()
	}()

	return &QuitResponse{}, nil
}
