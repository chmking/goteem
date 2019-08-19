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
	"github.com/chmking/horde/tsbuffer"
	"google.golang.org/grpc"
)

func New() *Manager {
	return &Manager{
		buffer: tsbuffer.New(time.Second * 5),
		sm:     &state.StateMachine{},
	}
}

type agentRegistry struct {
	Host string
	Port string

	Client private.AgentClient
	Status public.Status
}

type Manager struct {
	buffer *tsbuffer.Buffer
	agents []*agentRegistry
	sm     *state.StateMachine
	mtx    sync.Mutex

	cancel      context.CancelFunc
	tallyCancel context.CancelFunc
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

	tallyCtx, cancel := context.WithCancel(context.Background())
	m.tallyCancel = cancel

	go m.tally(tallyCtx)

	return &public.StartResponse{}, nil
}

func (m *Manager) Stop(ctx context.Context, req *public.StopRequest) (*public.StopResponse, error) {
	log.Println("Received request to stop")

	if m.tallyCancel != nil {
		m.tallyCancel()
		m.tallyCancel = nil
	}

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

func (m *Manager) Quit(ctx context.Context, req *public.QuitRequest) (*public.QuitResponse, error) {
	log.Println("Received request to quit")

	if m.tallyCancel != nil {
		m.tallyCancel()
		m.tallyCancel = nil
	}

	if err := m.sm.Quitting(); err != nil {
		return nil, err
	}

	m.mtx.Lock()
	for _, agent := range m.agents {
		agent.Client.Quit(ctx, &private.QuitRequest{})
	}
	m.mtx.Unlock()

	defer func() {
		m.cancel()
	}()

	return &public.QuitResponse{}, nil
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
	ctx, m.cancel = context.WithCancel(ctx)

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

func (m *Manager) tally(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				results := m.buffer.Collect()
				if results != nil {
					log.Printf("%+v", results)
				}
				// 	for _, result := range results {
				// 		log.Printf("%+v", result)
				// 	}
			}

			<-time.After(time.Second * 5)
		}
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
					for _, result := range resp.Results {
						m.buffer.Add(result)
					}
				}
				m.mtx.Unlock()

				<-time.After(time.Second)
			}
		}
	}()
}
