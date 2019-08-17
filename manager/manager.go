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

type Manager struct {
	agents []AgentRegistry
	sm     state.StateMachine
	mtx    sync.Mutex
}

type AgentRegistry struct {
	Host string
	Port string

	Client private.AgentClient
}

func (m *Manager) Start(ctx context.Context, req *public.StartRequest) (*public.StartResponse, error) {
	log.Println("Received private.StartRequest")
	if err := m.sm.Scaling(); err != nil {
		return nil, err
	}

	m.mtx.Lock()
	for _, agent := range m.agents {
		_, err := agent.Client.Scale(ctx, &private.ScaleRequest{})
		if err != nil {
			// agent should be quarantined
			log.Print(err)
		}
	}
	m.mtx.Unlock()

	return &public.StartResponse{}, nil
}

func (m *Manager) Stop(ctx context.Context, req *public.StopRequest) (*public.StopResponse, error) {
	log.Println("Received private.StopRequest")
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

func (m *Manager) Register(ctx context.Context, req *private.RegisterRequest) (*private.RegisterResponse, error) {
	host := req.Host + ":" + req.Port

	conn, err := grpc.Dial(host,
		grpc.WithBackoffMaxDelay(time.Second),
		grpc.WithInsecure())
	if err != nil {
		log.Print(err)
		return nil, err
	}

	registry := AgentRegistry{
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

func (m *Manager) ListenAndServePublic() {
	go func() {
		if err := m.sm.Idle(); err != nil {
			log.Fatal(err)
		}

		lis, err := net.Listen("tcp", ":8089")
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Listening for public connections on :8089")

		server := grpc.NewServer()
		public.RegisterManagerServer(server, m)
		log.Fatal(server.Serve(lis))
	}()
}

func (m *Manager) ListenAndServePrivate() {
	go func() {
		lis, err := net.Listen("tcp", ":5557")
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Listening for private connections on :5557")

		server := grpc.NewServer()
		private.RegisterManagerServer(server, m)
		log.Fatal(server.Serve(lis))
	}()
}

func (m *Manager) Healthcheck(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				m.mtx.Lock()
				for _, agent := range m.agents {
					_, err := agent.Client.Heartbeat(ctx, &private.HeartbeatRequest{})
					if err != nil {
						// Agent should be moved to 'Unhealthy' after some
						// number of failures.
						log.Print(err)
					}
				}
				m.mtx.Unlock()

				<-time.After(time.Second)
			}
		}
	}()
}
