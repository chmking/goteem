package manager

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/chmking/horde/helpers"
	"github.com/chmking/horde/manager/registry"
	"github.com/chmking/horde/manager/tsbuffer"
	"github.com/chmking/horde/protobuf/private"
	"github.com/chmking/horde/protobuf/public"
	"github.com/chmking/horde/state"
	"google.golang.org/grpc"
)

var ErrNoActiveAgents = errors.New("no active agents")

func New() *Manager {
	manager := &Manager{
		Registry:     registry.New(),
		buffer:       tsbuffer.New(time.Second * 5),
		StateMachine: &state.StateMachine{},
	}

	manager.Registry.RegisterCallback(manager.Rebalance)
	manager.Registry.BeginHealthcheck(context.Background())

	return manager
}

type Registry interface {
	Add(registry.Registration) error

	GetAll() []registry.Registration
	GetActive() []registry.Registration

	RegisterCallback(func())
	BeginHealthcheck(context.Context)
}

var _ = Registry(&registry.Registry{})

type StateMachine interface {
	State() public.Status

	Scaling() error
	Stopping() error
}

type agentRegistry struct {
	Host string
	Port string

	Client private.AgentClient
	Status public.Status
}

type WorkOrder struct {
	Id    string
	Users int
	Rate  float64
}

type Manager struct {
	Registry     Registry
	buffer       *tsbuffer.Buffer
	StateMachine StateMachine

	current WorkOrder

	cancel context.CancelFunc
}

func (m *Manager) State() public.Status {
	return m.StateMachine.State()
}

func (m *Manager) Start(count int, rate float64) error {
	if len(m.Registry.GetActive()) == 0 {
		return ErrNoActiveAgents
	}

	currentState := m.StateMachine.State()
	if err := m.StateMachine.Scaling(); err != nil {
		return err
	}

	if currentState == public.Status_STATUS_IDLE {
		m.current = WorkOrder{
			Id: helpers.MustUUID(),
		}
	}

	m.current.Users = count
	m.current.Rate = rate

	m.AssignWorkOrder(context.Background(), m.current)

	return nil
}

func (m *Manager) Stop() error {
	if err := m.StateMachine.Stopping(); err != nil {
		return err
	}

	all := m.Registry.GetAll()
	for _, agent := range all {
		_, err := agent.Client.Stop(context.Background(), &private.StopRequest{})
		if err != nil {
			log.Print(err)
		}
	}

	return nil
}

func (m *Manager) AssignWorkOrder(ctx context.Context, order WorkOrder) {
	active := m.Registry.GetActive()
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

type Registration struct {
	Id   string
	Host string
	Port string
}

func (m *Manager) Register(req Registration) error {
	address := req.Host + ":" + req.Port
	conn, err := grpc.Dial(address,
		grpc.WithBackoffMaxDelay(time.Second),
		grpc.WithInsecure())
	if err != nil {
		return err
	}

	regis := registry.Registration{
		Id:     req.Id,
		Client: private.NewAgentClient(conn),
	}

	m.Registry.Add(regis)

	return nil
}

func (m *Manager) Rebalance() {
	current := m.StateMachine.State()
	if current != public.Status_STATUS_RUNNING &&
		current != public.Status_STATUS_SCALING {
		return
	}

	m.AssignWorkOrder(context.Background(), m.current)
}
