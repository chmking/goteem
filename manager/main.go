package main

import (
	"context"
	"log"
	"net"

	"github.com/chmking/horde/protobuf/private"
	"github.com/chmking/horde/protobuf/public"
	"google.golang.org/grpc"
)

type Manager struct {
}

func (m *Manager) Start(ctx context.Context, req *public.StartRequest) (*public.StartResponse, error) {
	log.Println("Received a public.Start request")
	return &public.StartResponse{}, nil
}

func (m *Manager) Stop(ctx context.Context, req *public.StopRequest) (*public.StopResponse, error) {
	log.Println("Received a public.Stop request")
	return &public.StopResponse{}, nil
}

func (m *Manager) Heartbeat(ctx context.Context, req *private.HeartbeatRequest) (*private.HeartbeatResponse, error) {
	log.Println("Received a priavte.Heartbeat request")
	return &private.HeartbeatResponse{}, nil
}

func main() {

	manager := &Manager{}

	// Start the public endpoints
	go func() {
		lis, err := net.Listen("tcp", ":8089")
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Listening for public connections on :8089")

		server := grpc.NewServer()
		public.RegisterManagerServer(server, manager)
		log.Fatal(server.Serve(lis))
	}()

	lis, err := net.Listen("tcp", ":5557")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Listening for private connections on :5557")

	server := grpc.NewServer()
	private.RegisterManagerServer(server, manager)
	log.Fatal(server.Serve(lis))
}
