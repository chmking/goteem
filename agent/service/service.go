package service

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/chmking/horde"
	"github.com/chmking/horde/agent"
	"github.com/chmking/horde/protobuf/private"
	pb "github.com/chmking/horde/protobuf/private"
	"google.golang.org/grpc"
)

var _ = private.AgentServer(&Service{})

type Agent interface {
}

func New(config horde.Config) *Service {
	return &Service{
		agent: agent.New(config),
	}
}

type Service struct {
	port string

	agent  Agent
	cancel context.CancelFunc
}

func (s *Service) Healthcheck(
	ctx context.Context,
	req *pb.HealthcheckRequest) (*pb.HealthcheckResponse, error) {

	return &pb.HealthcheckResponse{}, nil
}

func (s *Service) Scale(
	ctx context.Context,
	req *pb.Orders) (*pb.ScaleResponse, error) {

	log.Println("Received request to scale")
	return &pb.ScaleResponse{}, nil
}

func (s *Service) Stop(
	ctx context.Context,
	req *pb.StopRequest) (*pb.StopResponse, error) {

	log.Println("Received request to stop")
	return &pb.StopResponse{}, nil
}

func (s *Service) Quit(
	ctx context.Context,
	req *pb.QuitRequest) (*pb.QuitResponse, error) {

	log.Println("Received request to quit")
	return &pb.QuitResponse{}, nil
}

func (s *Service) Listen(ctx context.Context) error {
	ctx, s.cancel = context.WithCancel(ctx)

	// TODO: Initialize agent

	errs := make(chan error, 1)

	s.listenAndServePrivate(errs)
	// s.dialManager(ctx, errs)

	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-errs:
			return err
		default:
			<-time.After(time.Second)
		}
	}
}

func (s *Service) listenAndServePrivate(errs chan<- error) {
	go func() {
		address := ":" + s.port
		lis, err := net.Listen("tcp", address)
		if err != nil {
			errs <- err
			return
		}

		log.Printf("Listening for private connection on %s\n", address)

		server := grpc.NewServer()
		pb.RegisterAgentServer(server, s)
		errs <- server.Serve(lis)
	}()
}
