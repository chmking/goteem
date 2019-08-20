package manager

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	"github.com/chmking/horde/helpers"
	"github.com/chmking/horde/protobuf/private"
	"github.com/chmking/horde/protobuf/public"
	"github.com/chmking/horde/state"
	"github.com/chmking/horde/tsbuffer"
	"google.golang.org/grpc"
)

func New() *Manager {
	manager := &Manager{
		registry: NewRegistry(),
		buffer:   tsbuffer.New(time.Second * 5),
		sm:       &state.StateMachine{},
	}

	manager.registry.RegisterCallback(manager.Rebalance)
	manager.registry.BeginHealthcheck(context.Background())

	return manager
}

type agentRegistry struct {
	Host string
	Port string

	Client private.AgentClient
	Status public.Status
}

type stateMachine interface {
	Idle() error
	Scaling() error
	Stopping() error
	Quitting() error
	State() public.Status
}

type WorkOrder struct {
	Id    string
	Users int
	Rate  float64
}

type Manager struct {
	registry *Registry
	buffer   *tsbuffer.Buffer
	agents   []*agentRegistry
	sm       stateMachine
	mtx      sync.Mutex

	current WorkOrder

	cancel      context.CancelFunc
	tallyCancel context.CancelFunc
}

func (m *Manager) State() public.Status {
	return m.sm.State()
}

func (m *Manager) Start(ctx context.Context, req *public.StartRequest) (*public.StartResponse, error) {
	log.Println("Receieved request to start")

	currentState := m.sm.State()

	if err := m.sm.Scaling(); err != nil {
		return nil, err
	}

	if currentState == public.Status_STATUS_IDLE {
		m.current = WorkOrder{
			Id: helpers.MustUUID(),
		}
	}

	m.current.Users = int(req.Users)
	m.current.Rate = req.Rate

	m.AssignWorkOrder(ctx, m.current)

	tallyCtx, cancel := context.WithCancel(context.Background())
	m.tallyCancel = cancel

	go m.tally(tallyCtx)

	return &public.StartResponse{}, nil
}

func (m *Manager) AssignWorkOrder(ctx context.Context, order WorkOrder) {
	active := m.registry.GetActive()
	activeLen := len(active)

	if order.Users < activeLen {
		activeLen = order.Users
	}

	allotment := order.Users / activeLen
	remainder := order.Users % activeLen
	increment := int64(float64(time.Second.Nanoseconds()) / order.Rate)

	for i := 0; i < activeLen; i++ {
		users := allotment
		if i < remainder {
			users = users + 1
		}
		rate := int64(activeLen) * increment
		wait := int64(i) * increment

		req := &private.ScaleRequest{
			Users: int32(users),
			Rate:  rate,
			Wait:  wait,
		}

		agent := active[0]
		_, err := agent.Client.Scale(ctx, req)
		if err != nil {
			// TODO: Quarantine agent
		}
	}
}

func (m *Manager) Stop(ctx context.Context, req *public.StopRequest) (*public.StopResponse, error) {
	log.Println("Received request to stop")

	if err := m.sm.Stopping(); err != nil {
		return nil, err
	}

	if m.tallyCancel != nil {
		m.tallyCancel()
		m.tallyCancel = nil
	}

	all := m.registry.GetAll()
	for _, agent := range all {
		_, err := agent.Client.Stop(context.Background(), &private.StopRequest{})
		if err != nil {
			log.Print(err)
		}
	}

	return &public.StopResponse{}, nil
}

func (m *Manager) Quit(ctx context.Context, req *public.QuitRequest) (*public.QuitResponse, error) {
	log.Println("Received request to quit")

	if err := m.sm.Quitting(); err != nil {
		return nil, err
	}

	if m.tallyCancel != nil {
		m.tallyCancel()
		m.tallyCancel = nil
	}

	all := m.registry.GetAll()
	for _, agent := range all {
		_, err := agent.Client.Quit(context.Background(), &private.QuitRequest{})
		if err != nil {
			log.Print(err)
		}
	}

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
	log.Printf("Receivied regitration request: %+v", req)

	address := req.Host + ":" + req.Port

	conn, err := grpc.Dial(address,
		grpc.WithBackoffMaxDelay(time.Second),
		grpc.WithInsecure())
	if err != nil {
		log.Print(err)
		return nil, err
	}

	regis := Registration{
		Id:     req.Id,
		Client: private.NewAgentClient(conn),
	}

	log.Printf("Adding agent to registry: %+v\n", regis)
	m.registry.Add(regis)

	return &private.RegisterResponse{}, nil
}

func (m *Manager) Rebalance() {
	current := m.sm.State()
	if current != public.Status_STATUS_RUNNING &&
		current != public.Status_STATUS_SCALING {
		return
	}

	m.AssignWorkOrder(context.Background(), m.current)
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
