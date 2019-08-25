package service

import (
	"context"
	"net"

	"github.com/chmking/horde/logger/log"
	"github.com/chmking/horde/manager"
	"github.com/chmking/horde/protobuf/private"
	"github.com/chmking/horde/protobuf/public"
	"google.golang.org/grpc"
)

var _ = public.ManagerServer(&Service{})
var _ = private.ManagerServer(&Service{})

type Manager interface {
	Start(count int, rate float64) error
	Stop() error

	Register(id, address string) error
}

var _ = Manager(&manager.Manager{})

func New() *Service {
	return &Service{
		manager: manager.New(),
	}
}

type Service struct {
	manager Manager
	cancel  context.CancelFunc
}

func (s *Service) Start(
	ctx context.Context,
	req *public.StartRequest) (*public.StartResponse, error) {

	log.Info().Msg("Receieved request to start")
	if err := s.manager.Start(int(req.Users), req.Rate); err != nil {
		return nil, err
	}

	return &public.StartResponse{}, nil
}

func (s *Service) Status(
	ctx context.Context,
	req *public.StatusRequest) (*public.StatusResponse, error) {

	log.Info().Msg("Receivied request for status")
	return &public.StatusResponse{}, nil
}

func (s *Service) Stop(
	ctx context.Context,
	req *public.StopRequest) (*public.StopResponse, error) {

	log.Info().Msg("Received request to stop")
	if err := s.manager.Stop(); err != nil {
		return nil, err
	}

	return &public.StopResponse{}, nil
}

func (s *Service) Quit(
	ctx context.Context,
	req *public.QuitRequest) (*public.QuitResponse, error) {

	log.Info().Msg("Received request to quit")
	if err := s.manager.Stop(); err != nil {
		return nil, err
	}

	defer func() {
		if s.cancel != nil {
			s.cancel()
		}
	}()

	return &public.QuitResponse{}, nil
}

func (s *Service) Register(
	ctx context.Context,
	req *private.RegisterRequest) (*private.RegisterResponse, error) {

	log.Info().Msg("Receivied request to register")
	if err := s.manager.Register(req.Id, req.Host+":"+req.Port); err != nil {
		return nil, err
	}

	return &private.RegisterResponse{}, nil
}

func (s *Service) Listen(ctx context.Context) error {
	ctx, s.cancel = context.WithCancel(ctx)

	errs := make(chan error, 1)
	s.listenAndServePublic(errs)
	s.listenAndServePrivate(errs)

	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-errs:
			return err
		}
	}
}

func (s *Service) listenAndServePublic(errs chan<- error) {
	go func() {
		lis, err := net.Listen("tcp", ":8089")
		if err != nil {
			errs <- err
			return
		}

		log.Info().Msg("Listening for public connections on :8089")

		server := grpc.NewServer()
		public.RegisterManagerServer(server, s)
		errs <- server.Serve(lis)
	}()
}

func (m *Service) listenAndServePrivate(errs chan<- error) {
	go func() {
		lis, err := net.Listen("tcp", ":5557")
		if err != nil {
			errs <- err
			return
		}

		log.Info().Msg("Listening for private connections on :5557")

		server := grpc.NewServer()
		private.RegisterManagerServer(server, m)
		errs <- server.Serve(lis)
	}()
}
