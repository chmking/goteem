package agent

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/chmking/horde"
	"github.com/chmking/horde/agent/recorder"
	"github.com/chmking/horde/agent/session"
	"github.com/chmking/horde/helpers"
	pb "github.com/chmking/horde/protobuf/private"
	"github.com/chmking/horde/protobuf/public"
	"github.com/chmking/horde/state"
	grpc "google.golang.org/grpc"
)

type stateMachine interface {
	Idle() error
	Running() error
	Scaling() error
	Stopping() error
	State() public.Status
}

func New(config horde.Config) *Agent {
	uuid := helpers.MustUUID()
	hostname := helpers.MustHostname()

	return &Agent{
		id:       hostname + "_" + uuid,
		hostname: hostname,
		host:     agentHost,
		port:     agentPort,

		config:   config,
		recorder: recorder.New(),
		session:  &session.Session{},
		sm:       &state.StateMachine{},
	}
}

type Agent struct {
	id       string
	hostname string
	host     string
	port     string

	config   horde.Config
	recorder *recorder.Recorder
	session  *session.Session
	sm       stateMachine
	server   *grpc.Server
	mtx      sync.Mutex

	cancel context.CancelFunc
}

func (a *Agent) Scale(count int32, rate, wait int64) error {
	if err := a.sm.Scaling(); err != nil {
		return err
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
	a.session.Scale(ctx, order, a.onScaled)

	return nil
}

func (a *Agent) Stop() error {
	if err := a.sm.Stopping(); err != nil {
		return err
	}

	a.session.Stop(a.onStopped)

	return nil
}

func (a *Agent) onScaled() {
	a.sm.Running()
}

func (a *Agent) onStopped() {
	a.sm.Idle()
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
