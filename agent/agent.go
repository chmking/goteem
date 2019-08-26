package agent

import (
	"context"
	"sync"

	"github.com/chmking/horde"
	"github.com/chmking/horde/agent/recorder"
	"github.com/chmking/horde/agent/session"
	"github.com/chmking/horde/eventloop"
	"github.com/chmking/horde/logger/log"
	"github.com/chmking/horde/state"
	grpc "google.golang.org/grpc"
)

func New(config horde.Config) *Agent {

	agent := &Agent{
		StateMachine: &state.StateMachine{},
		Session:      &session.Session{},

		config:   config,
		recorder: recorder.New(),
		events:   eventloop.New(),
	}

	agent.StateMachine.Idle()

	return agent
}

type StateMachine interface {
	Idle() error
	Running() error
	Scaling() error
	Stopping() error
}

var _ = StateMachine(&state.StateMachine{})

type Session interface {
	Scale(context.Context, session.ScaleOrder, session.Callback)
	Stop(session.Callback)
}

var _ = Session(&session.Session{})

type Agent struct {
	StateMachine StateMachine
	Session      Session

	config   horde.Config
	recorder *recorder.Recorder
	events   *eventloop.EventLoop
	server   *grpc.Server
	mtx      sync.Mutex

	cancel context.CancelFunc
}

type Orders struct {
	Id string

	Count int32
	Rate  int64
	Wait  int64
}

func (a *Agent) Scale(orders Orders) (err error) {
	a.events.Append(func() {
		if err = a.StateMachine.Scaling(); err != nil {
			return
		}

		sessionOrders := session.ScaleOrder{
			Count: orders.Count,
			Rate:  orders.Rate,
			Wait:  orders.Wait,
			Work: session.Work{
				Tasks:   a.config.Tasks,
				WaitMin: a.config.WaitMin,
				WaitMax: a.config.WaitMax,
			},
		}

		log.Info().Msgf("Requesting Scale with ScaleOrder: %+v", sessionOrders)

		ctx := horde.WithRecorder(context.Background(), a.recorder)
		a.Session.Scale(ctx, sessionOrders, a.onScaled)
	})

	return
}

func (a *Agent) Stop() (err error) {
	a.events.Append(func() {
		if err = a.StateMachine.Stopping(); err != nil {
			return
		}

		a.Session.Stop(a.onStopped)
	})

	return
}

func (a *Agent) onScaled() {
	a.events.Append(func() {
		a.StateMachine.Running()
	})
}

func (a *Agent) onStopped() {
	a.events.Append(func() {
		a.StateMachine.Idle()
	})
}
