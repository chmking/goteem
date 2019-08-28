package service

import (
	"context"
	"net"
	"time"

	"github.com/chmking/horde"
	"github.com/chmking/horde/agent"
	"github.com/chmking/horde/helpers"
	"github.com/chmking/horde/logger/log"
	"github.com/chmking/horde/protobuf/private"
	pb "github.com/chmking/horde/protobuf/private"
	"google.golang.org/grpc"
)

var _ = private.AgentServer(&Service{})

type Agent interface {
	Status() agent.Status
	Scale(order agent.Orders) error
	Stop() error
}

func New(config horde.Config) *Service {
	uuid := helpers.MustUUID()
	hostname := helpers.MustHostname()

	return &Service{
		id:       hostname + "_" + uuid,
		hostname: hostname,
		host:     agentHost,
		port:     agentPort,

		Agent: agent.New(config),
	}
}

type Service struct {
	id       string
	hostname string
	host     string
	port     string

	Agent  Agent
	cancel context.CancelFunc
}

func (s *Service) Healthcheck(
	ctx context.Context,
	req *pb.HealthcheckRequest) (*pb.HealthcheckResponse, error) {

	status := s.Agent.Status()
	resp := &pb.HealthcheckResponse{
		State: status.State,
		Count: int32(status.Count),
	}

	return resp, nil
}

func (s *Service) Scale(
	ctx context.Context,
	req *pb.ScaleRequest) (*pb.ScaleResponse, error) {

	log.Info().Msg("Received request to scale")

	orders := agent.Orders{
		Id: req.Orders.Id,

		Count: req.Orders.Count,
		Rate:  req.Orders.Rate,
		Wait:  req.Orders.Wait,
	}

	if err := s.Agent.Scale(orders); err != nil {
		return &pb.ScaleResponse{}, err
	}

	return &pb.ScaleResponse{}, nil
}

func (s *Service) Stop(
	ctx context.Context,
	req *pb.StopRequest) (*pb.StopResponse, error) {

	log.Info().Msg("Received request to stop")
	if err := s.Agent.Stop(); err != nil {
		return &pb.StopResponse{}, err
	}

	return &pb.StopResponse{}, nil
}

func (s *Service) Quit(
	ctx context.Context,
	req *pb.QuitRequest) (*pb.QuitResponse, error) {

	log.Info().Msg("Received request to quit")
	if err := s.Agent.Stop(); err != nil {
		return &pb.QuitResponse{}, err
	}

	return &pb.QuitResponse{}, nil
}

func (s *Service) Listen(ctx context.Context) error {
	ctx, s.cancel = context.WithCancel(ctx)

	errs := make(chan error, 1)
	s.listenAndServePrivate(errs)
	s.dialManager(errs)

	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-errs:
			return err
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

		log.Info().Msgf("Listening for private connection on %s", address)

		server := grpc.NewServer()
		pb.RegisterAgentServer(server, s)
		errs <- server.Serve(lis)
	}()
}

func (s *Service) dialManager(errs chan<- error) {
	go func() {
		address := managerHost + ":" + managerPort
		conn, err := grpc.Dial(address,
			grpc.WithBackoffMaxDelay(time.Second),
			grpc.WithInsecure(),
			grpc.WithBlock())
		if err != nil {
			errs <- err
			return
		}

		client := pb.NewManagerClient(conn)

		req := &pb.RegisterRequest{
			Id:       s.id,
			Hostname: s.hostname,
			Host:     s.host,
			Port:     s.port,
		}

		_, err = client.Register(context.Background(), req)
		if err != nil {
			errs <- err
			return
		}
	}()
}
