package main

import (
	"context"
	"log"
	"net"

	"github.com/chmking/horde/protobuf/public"
	"google.golang.org/grpc"
)

type Manager struct {
}

func (m *Manager) Start(ctx context.Context, req *public.StartRequest) (*public.StartResponse, error) {
	log.Println("Received a Start request")
	return &public.StartResponse{}, nil
}

func (m *Manager) Stop(ctx context.Context, req *public.StopRequest) (*public.StopResponse, error) {
	log.Println("Received a Stop request")
	return &public.StopResponse{}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8089")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Listening to requests on :8089")

	manager := &Manager{}
	server := grpc.NewServer()
	public.RegisterManagerServer(server, manager)
	log.Fatal(server.Serve(lis))
}
