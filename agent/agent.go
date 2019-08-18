package agent

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	"github.com/chmking/horde"
	pb "github.com/chmking/horde/protobuf/private"
	"github.com/chmking/horde/protobuf/public"
	"github.com/chmking/horde/recorder"
	"github.com/chmking/horde/session"
	sess "github.com/chmking/horde/session"
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
	return &Agent{
		config:   config,
		recorder: recorder.New(),
		session:  &sess.Session{},
		sm:       &state.StateMachine{},
	}
}

type Agent struct {
	config   horde.Config
	recorder *recorder.Recorder
	session  *session.Session
	sm       stateMachine
	server   *grpc.Server
	mtx      sync.Mutex
}

func (a *Agent) Listen(ctx context.Context) error {
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
		lis, err := net.Listen("tcp", ":5558")
		if err != nil {
			errs <- err
			return
		}

		log.Println("Listening for private connection on :5558")

		a.server = grpc.NewServer()
		pb.RegisterAgentServer(a.server, a)
		errs <- a.server.Serve(lis)
	}()
}

func (a *Agent) dialManager(ctx context.Context, errs chan<- error) {
	go func() {
		conn, err := grpc.Dial("127.0.0.1:5557",
			grpc.WithBackoffMaxDelay(time.Second),
			grpc.WithInsecure(),
			grpc.WithBlock())
		if err != nil {
			errs <- err
			return
		}

		client := pb.NewManagerClient(conn)

		req := &pb.RegisterRequest{
			Host: "127.0.0.1",
			Port: "5558",
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

	order := sess.ScaleOrder{
		Count: req.Users,
		Rate:  req.Rate,
		Wait:  req.Wait,
		Work: sess.Work{
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
