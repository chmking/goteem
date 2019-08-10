package horde

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
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	a.server = grpc.NewServer()
	RegisterAgentServer(a.server, a)
	return a.server.Serve(lis)
}

func (a *Agent) Start(ctx context.Context, req *StartRequest) (*StartResponse, error) {
	return &StartResponse{}, nil
}

func (a *Agent) Quit(ctx context.Context, req *QuitRequest) (*QuitResponse, error) {
	defer func() {
		a.server.Stop()
	}()

	return &QuitResponse{}, nil
}
