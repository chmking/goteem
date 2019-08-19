package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chmking/horde"
	"github.com/chmking/horde/agent"
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
		log.Fatal(http.ListenAndServe(":6060", nil))
	}()

	simple := &horde.Task{
		Name: "hello_world",
		Func: func(ctx context.Context) {
			recorder := horde.RecorderFrom(ctx)

			start := time.Now()
			elapsed := time.Since(start).Nanoseconds() / 1e6

			recorder.Success("GET", "http://google.com", elapsed)
		},
	}

	config := horde.Config{
		Tasks:   []*horde.Task{simple},
		WaitMin: 100,
		WaitMax: 150,
	}

	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	agent := agent.New(config)

	go func() {
		if err := agent.Listen(ctx); err != nil {
			log.Print(err)
		}
		c <- syscall.SIGQUIT
	}()

	s := <-c
	fmt.Println("Got signal:", s)
	cancel()
}
