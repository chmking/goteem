package horde_test

import (
	"context"
	"net/http"
	"time"

	"github.com/chmking/horde"
)

// Example represents a custom test to be performed.
func Example(ctx context.Context) {
	// Recorder is embedded in the context
	recorder := horde.RecorderFrom(ctx)

	// Measure the request
	start := time.Now()
	_, err := http.Get("https://google.com")
	elapsed := time.Since(start).Nanoseconds() / time.Millisecond.Nanoseconds()

	// Record error
	if err != nil {
		recorder.Error("HTTP", "http://google.com", elapsed, err)
	}

	// Record success
	recorder.Success("HTTP", "http://google.com", elapsed)
}

func AgentMain() {
	// Define a simple HTTP task.
	task := &horde.Task{
		Name:   "example",
		Func:   Example,
		Weight: 1,
	}

	// Define a simple behavior.
	behavior := &horde.Behavior{
		Tasks: []*horde.Task{task},
	}

	// Create an agent.
	agent := horde.Agent{
		Behavior: behavior,
	}

	// Start listening.
	agent.Listen(":5557")
}

func main() {
	// Start the agent in a separate routine
	// for example purposes.
	go AgentMain()
}
