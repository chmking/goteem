package agent

import (
	"context"
	"log"
	"net"
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

func (a *Agent) Listen(ctx context.Context) error {
	ctx, a.cancel = context.WithCancel(ctx)

	if err := a.sm.Idle(); err != nil {
		return err
	}

	errs := make(chan error, 1)

	a.listenAndServePrivate(errs)
	a.dialManager(ctx, errs)

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

func (a *Agent) listenAndServePrivate(errs chan<- error) {
	go func() {
		address := ":" + a.port
		lis, err := net.Listen("tcp", address)
		if err != nil {
			errs <- err
			return
		}

		log.Printf("Listening for private connection on %s\n", address)

		a.server = grpc.NewServer()
		pb.RegisterAgentServer(a.server, a)
		errs <- a.server.Serve(lis)
	}()
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

func (a *Agent) Scale(ctx context.Context, req *pb.ScaleRequest) (*pb.ScaleResponse, error) {
	log.Println("Received request to scale")

	if err := a.sm.Scaling(); err != nil {
		return nil, err
	}

	order := session.ScaleOrder{
		Count: req.Users,
		Rate:  req.Rate,
		Wait:  req.Wait,
		Work: session.Work{
			Tasks:   a.config.Tasks,
			WaitMin: a.config.WaitMin,
			WaitMax: a.config.WaitMax,
		},
	}

	log.Printf("Requesting Scale with ScaleOrder: %+v", order)

	rctx := horde.WithRecorder(context.Background(), a.recorder)
	a.session.Scale(rctx, order, a.onScaled)

	return &pb.ScaleResponse{}, nil
}

func (a *Agent) Stop(ctx context.Context, req *pb.StopRequest) (*pb.StopResponse, error) {
	log.Println("Received request to stop work")

	if err := a.sm.Stopping(); err != nil {
		return nil, err
	}

	a.session.Stop(a.onStopped)

	return &pb.StopResponse{}, nil
}

func (a *Agent) Quit(ctx context.Context, req *pb.QuitRequest) (*pb.QuitResponse, error) {
	log.Println("Received private.QuitRequest")

	if err := a.sm.Stopping(); err != nil {
		return nil, err
	}

	defer a.cancel()

	return &pb.QuitResponse{}, nil
}

func (a *Agent) Heartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	results := a.recorder.Results()

	resp := &pb.HeartbeatResponse{
		Status:  a.sm.State(),
		Results: results,
	}

	return resp, nil
}

func (a *Agent) onScaled() {
	a.sm.Running()
}

func (a *Agent) onStopped() {
	a.sm.Idle()
}
