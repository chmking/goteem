package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"time"

	"github.com/chmking/horde"
	"github.com/chmking/horde/agent"
)

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

	go func() {
		agent := agent.New(config)
		log.Fatal(agent.Listen(ctx))
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	s := <-c
	fmt.Println("Got signal:", s)
	cancel()
}
