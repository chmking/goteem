package main

import (
	"context"
	"flag"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chmking/horde"
	"github.com/chmking/horde/agent"
	"github.com/chmking/horde/agent/service"
	"github.com/chmking/horde/logger/log"
)

func init() {
	// Register the agent flags
	agent.Flags()

	// Add additional flags

	// Parse all flags
	flag.Parse()
}

func main() {
	// task := &horde.Task{
	// 	Name: "example",
	// 	Func: func(ctx context.Context) {
	// 		recorder := horde.RecorderFrom(ctx)

	// 		start := time.Now()
	// 		_, err := http.Get("https://google.com")
	// 		elapsed := time.Since(start).Nanoseconds()

	// 		if err != nil {
	// 			recorder.Error("GET", "http://google.com", elapsed, err)
	// 		}

	// 		recorder.Success("GET", "http://google.com", elapsed)
	// 	},
	// }

	go func() {
		err := http.ListenAndServe(":6060", nil)
		log.Fatal().Err(err).Msg("pprof quit unexpectedly")
	}()

	simple := &horde.Task{
		Name: "hello_world",
		Func: func(ctx context.Context) {
			recorder := horde.RecorderFrom(ctx)

			start := time.Now()
			log.Info().Msg("Sample")
			elapsed := time.Since(start).Nanoseconds() / 1e6

			recorder.Success("GET", "http://google.com", elapsed)
		},
	}

	config := horde.Config{
		Tasks:   []*horde.Task{simple},
		WaitMin: 1000,
		WaitMax: 1000,
	}

	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	svc := service.New(config)

	go func() {
		if err := svc.Listen(ctx); err != nil {
			log.Error().Err(err).Msg("agent quit unexpectedly")
		}
		c <- syscall.SIGQUIT
	}()

	s := <-c
	log.Info().Msg("Got signal: " + s.String())
	cancel()
}
