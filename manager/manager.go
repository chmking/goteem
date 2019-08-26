package manager

import (
	"context"
	"errors"
	"time"

	"github.com/chmking/horde/eventloop"
	"github.com/chmking/horde/helpers"
	"github.com/chmking/horde/logger/log"
	"github.com/chmking/horde/manager/registry"
	"github.com/chmking/horde/manager/tsbuffer"
	"github.com/chmking/horde/protobuf/private"
	"github.com/chmking/horde/protobuf/public"
	"github.com/chmking/horde/state"
	"google.golang.org/grpc"
)

var ErrNoActiveAgents = errors.New("no active agents")

type Registry interface {
	Add(registry.Registration) error
	Quarantine(string) error

	GetAll() []registry.Registration
	GetActive() []registry.Registration

	RegisterCallback(func())
	BeginHealthcheck(context.Context)
}

var _ = Registry(&registry.Registry{})

type StateMachine interface {
	State() public.Status

	Idle() error
	Scaling() error
	Stopping() error
}

var _ = StateMachine(&state.StateMachine{})

func New() *Manager {
	manager := &Manager{
		Registry:     registry.New(),
		StateMachine: &state.StateMachine{},

		buffer: tsbuffer.New(time.Second * 5),
		events: eventloop.New(),
	}

	manager.Registry.RegisterCallback(manager.OnRebalance)
	manager.Registry.BeginHealthcheck(context.Background())

	manager.StateMachine.Idle()

	return manager
}

type Orders struct {
	Id string

	Count int
	Rate  float64
}

type Manager struct {
	Registry     Registry
	StateMachine StateMachine

	buffer *tsbuffer.Buffer
	events *eventloop.EventLoop

	orders Orders
	cancel context.CancelFunc
}

func (m *Manager) State() (state public.Status) {
	m.events.Append(func() {
		state = m.StateMachine.State()
	})

	return
}

func (m *Manager) Start(count int, rate float64) (err error) {
	m.events.Append(func() {
		// Check for active agents
		log.Info().Msg("Getting active agents")
		if len(m.Registry.GetActive()) == 0 {
			err = ErrNoActiveAgents
			return
		}

		// Change state
		prevState := m.StateMachine.State()
		if err = m.StateMachine.Scaling(); err != nil {
			return
		}

		// Create new work when idle
		if prevState == public.Status_STATUS_IDLE {
			m.orders.Id = helpers.MustUUID()
		}

		// Update the work
		m.orders.Count = count
		m.orders.Rate = rate

		// Assign the work
		m.assignOrders()
	})

	return
}

func (m *Manager) Stop() (err error) {
	m.events.Append(func() {
		// Change state
		if err = m.StateMachine.Stopping(); err != nil {
			return
		}

		// Request agents stop
		all := m.Registry.GetAll()
		for _, agent := range all {
			if _, err := agent.Client.Stop(context.Background(), &private.StopRequest{}); err != nil {
				log.Error().Err(err).Msg("agent errored on stop request")
				m.Registry.Quarantine(agent.Id)
			}
		}
	})

	return
}

func (m *Manager) Register(id, address string) (err error) {
	m.events.Append(func() {
		var conn *grpc.ClientConn
		conn, err = grpc.Dial(address,
			grpc.WithBackoffMaxDelay(time.Second),
			grpc.WithInsecure())
		if err != nil {
			return
		}

		regis := registry.Registration{
			Id:     id,
			Client: private.NewAgentClient(conn),
		}

		log.Info().Msgf("Adding registry for: %+v", regis)
		m.Registry.Add(regis)
	})

	return
}

func (m *Manager) OnRebalance() {
	m.events.Append(func() {
		current := m.StateMachine.State()
		if !(current == public.Status_STATUS_RUNNING || current == public.Status_STATUS_SCALING) {
			return
		}

		m.assignOrders()
	})
}

func (m *Manager) assignOrders() {
	if m.orders.Count == 0 || m.orders.Rate == 0 {
		return
	}

	active := m.Registry.GetActive()
	activeLen := len(active)

	if m.orders.Count < activeLen {
		activeLen = m.orders.Count
	}

	allotment := m.orders.Count / activeLen
	remainder := m.orders.Count % activeLen
	increment := int64(float64(time.Second.Nanoseconds()) / m.orders.Rate)

	for i := 0; i < activeLen; i++ {
		count := allotment
		if i < remainder {
			count = count + 1
		}
		rate := int64(activeLen) * increment
		wait := int64(i) * increment

		req := &private.ScaleRequest{
			Orders: &private.Orders{
				Id: m.orders.Id,

				Count: int32(count),
				Rate:  rate,
				Wait:  wait,
			},
		}

		agent := active[0]
		_, err := agent.Client.Scale(context.Background(), req)
		if err != nil {
			log.Error().Err(err).Msg("agent errored on scale request")
			m.Registry.Quarantine(agent.Id)
		}
	}
}
