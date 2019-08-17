package main

import (
	"context"
	"log"

	"github.com/chmking/horde"
	"github.com/chmking/horde/agent"
)

func main() {
	// 	task := &horde.Task{
	// 		Name: "example",
	// 		Func: func(ctx context.Context) {
	// 			recorder := horde.RecorderFrom(ctx)
	//
	// 			start := time.Now()
	// 			_, err := http.Get("https://google.com")
	// 			elapsed := time.Since(start).Nanoseconds()
	//
	// 			if err != nil {
	// 				recorder.Error("GET", "http://google.com", elapsed, err)
	// 			}
	//
	// 			recorder.Success("GET", "http://google.com", elapsed)
	// 		},
	// 	}

	simple := &horde.Task{
		Name: "hello_world",
		Func: func(ctx context.Context) {
			log.Println("Hello World!")
		},
	}

	config := horde.Config{
		Tasks:   []*horde.Task{simple},
		WaitMin: 1000,
		WaitMax: 1500,
	}

	ctx, _ := context.WithCancel(context.Background())

	agent := agent.New(config)
	log.Fatal(agent.Listen(ctx))
}
