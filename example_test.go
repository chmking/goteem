package horde_test

import (
	"context"
	"net/http"
	"time"

	"github.com/chmking/horde"
	"github.com/chmking/horde/agent"
)

func AgentMain() {
	task := &horde.Task{
		Name: "example",
		Func: func(ctx context.Context) {
			recorder := horde.RecorderFrom(ctx)

			start := time.Now()
			_, err := http.Get("https://google.com")
			elapsed := time.Since(start).Nanoseconds()

			if err != nil {
				recorder.Error("GET", "http://google.com", elapsed, err)
			}

			recorder.Success("GET", "http://google.com", elapsed)
		},
	}

	config := horde.Config{
		Tasks:   []*horde.Task{task},
		WaitMin: 1000,
		WaitMax: 1500,
	}

	ctx := context.Background()

	agent := agent.New(config)
	agent.Listen(ctx)
}

func main() {
	go AgentMain()
}
