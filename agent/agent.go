package agent

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/chmking/horde"
	"github.com/chmking/horde/agent/recorder"
	"github.com/chmking/horde/agent/session"
	"github.com/chmking/horde/eventloop"
	"github.com/chmking/horde/helpers"
	pb "github.com/chmking/horde/protobuf/private"
	"github.com/chmking/horde/state"
	grpc "google.golang.org/grpc"
)

func New(config horde.Config) *Agent {
	uuid := helpers.MustUUID()
	hostname := helpers.MustHostname()

	return &Agent{
		StateMachine: &state.StateMachine{},
		Session:      &session.Session{},

		id:       hostname + "_" + uuid,
		hostname: hostname,
		host:     agentHost,
		port:     agentPort,

		config:   config,
		recorder: recorder.New(),
		events:   eventloop.New(),
	}
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

	id       string
	hostname string
	host     string
	port     string

	config   horde.Config
	recorder *recorder.Recorder
	events   *eventloop.EventLoop
	server   *grpc.Server
	mtx      sync.Mutex

	cancel context.CancelFunc
}

func (a *Agent) Scale(count int32, rate, wait int64) (err error) {
	a.events.Append(func() {
		if err = a.StateMachine.Scaling(); err != nil {
			return
		}

		order := session.ScaleOrder{
			Count: count,
			Rate:  rate,
			Wait:  wait,
			Work: session.Work{
				Tasks:   a.config.Tasks,
				WaitMin: a.config.WaitMin,
				WaitMax: a.config.WaitMax,
			},
		}

		log.Printf("Requesting Scale with ScaleOrder: %+v", order)

		ctx := horde.WithRecorder(context.Background(), a.recorder)
		a.Session.Scale(ctx, order, a.onScaled)
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

func (a *Agent) dialManager(ctx context.Context, errs chan<- error) {
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
			Id:       a.id,
			Hostname: a.hostname,
			Host:     a.host,
			Port:     a.port,
		}

		_, err = client.Register(ctx, req)
		if err != nil {
			errs <- err
			return
		}
	}()
}
