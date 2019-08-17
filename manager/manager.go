package manager

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	"github.com/chmking/horde/protobuf/private"
	"github.com/chmking/horde/protobuf/public"
	"github.com/chmking/horde/state"
	"google.golang.org/grpc"
)

func New() *Manager {
	return &Manager{
		sm: &state.StateMachine{},
	}
}

type stateMachine interface {
	Idle() error
	Scaling() error
	Stopping() error
	State() public.Status
}

type agentRegistry struct {
	Host string
	Port string

	Client private.AgentClient
	Status public.Status
}

type Manager struct {
	agents []*agentRegistry
	sm     stateMachine
	mtx    sync.Mutex
}

func (m *Manager) State() public.Status {
	return m.sm.State()
}

func (m *Manager) Start(ctx context.Context, req *public.StartRequest) (*public.StartResponse, error) {
	log.Println("Receieved request to start")

	if err := m.sm.Scaling(); err != nil {
		return nil, err
	}

	m.mtx.Lock()
	for _, agent := range m.agents {
		scaleReq := &private.ScaleRequest{
			Users: req.Users,
			Rate:  req.Rate,
		}

		_, err := agent.Client.Scale(ctx, scaleReq)
		if err != nil {
			// agent should be quarantined
			log.Print(err)
		}
	}
	m.mtx.Unlock()

	return &public.StartResponse{}, nil
}

func (m *Manager) Stop(ctx context.Context, req *public.StopRequest) (*public.StopResponse, error) {
	log.Println("Received request to stop")

	if err := m.sm.Stopping(); err != nil {
		return nil, err
	}

	m.mtx.Lock()
	for _, agent := range m.agents {
		_, err := agent.Client.Stop(ctx, &private.StopRequest{})
		if err != nil {
			// agent should be quarantined
			log.Print(err)
		}
	}
	m.mtx.Unlock()

	return &public.StopResponse{}, nil
}

func (m *Manager) Status(ctx context.Context, req *public.StatusRequest) (*public.StatusResponse, error) {
	resp := &public.StatusResponse{
		Status: m.sm.State(),
	}

	m.mtx.Lock()
	for _, agent := range m.agents {
		resp.Agents = append(resp.Agents, &public.AgentStatus{Status: agent.Status})
	}
	m.mtx.Unlock()

	return resp, nil
}

func (m *Manager) Register(ctx context.Context, req *private.RegisterRequest) (*private.RegisterResponse, error) {
	host := req.Host + ":" + req.Port

	conn, err := grpc.Dial(host,
		grpc.WithBackoffMaxDelay(time.Second),
		grpc.WithInsecure())
	if err != nil {
		log.Print(err)
		return nil, err
	}

	registry := &agentRegistry{
		Host:   req.Host,
		Port:   req.Port,
		Client: private.NewAgentClient(conn),
	}

	m.mtx.Lock()
	log.Printf("Adding agent to registry: %+v\n", registry)
	m.agents = append(m.agents, registry)
	m.mtx.Unlock()

	return &private.RegisterResponse{}, nil
}

func (m *Manager) Listen(ctx context.Context) error {
	if err := m.sm.Idle(); err != nil {
		return err
	}

	errs := make(chan error, 1)

	m.listenAndServePublic(errs)
	m.listenAndServePrivate(errs)
	m.healthcheck(ctx)

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

func (m *Manager) listenAndServePublic(errs chan<- error) {
	go func() {
		lis, err := net.Listen("tcp", ":8089")
		if err != nil {
			errs <- err
			return
		}

		log.Println("Listening for public connections on :8089")

		server := grpc.NewServer()
		public.RegisterManagerServer(server, m)
		errs <- server.Serve(lis)
	}()
}

func (m *Manager) listenAndServePrivate(errs chan<- error) {
	go func() {
		lis, err := net.Listen("tcp", ":5557")
		if err != nil {
			errs <- err
			return
		}

		log.Println("Listening for private connections on :5557")

		server := grpc.NewServer()
		private.RegisterManagerServer(server, m)
		errs <- server.Serve(lis)
	}()
}

func (m *Manager) healthcheck(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				m.mtx.Lock()
				for _, agent := range m.agents {
					resp, err := agent.Client.Heartbeat(ctx, &private.HeartbeatRequest{})
					if err != nil {
						// Agent should be moved to 'Unhealthy' after some
						// number of failures.
						log.Print(err)
						continue
					}

					agent.Status = resp.Status
				}
				m.mtx.Unlock()

				<-time.After(time.Second)
			}
		}
	}()
}
