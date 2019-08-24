package eventloop

import "context"

// Event represents an event in the loop
type Event func()

// New constructs a new EventLoop
func New() *EventLoop {
	e := &EventLoop{
		events: make(chan Event, 100),
	}
	go e.loop()

	return e
}

// EventLoop is a singles threaded eevnt loop
type EventLoop struct {
	events chan Event
}

// Append appends an event to the loop
func (e *EventLoop) Append(event Event) {
	ctx, cancel := context.WithCancel(context.Background())

	e.events <- func() {
		event()
		cancel()
	}

	<-ctx.Done()
}

func (e *EventLoop) loop() {
	for {
		event := <-e.events
		event()
	}
}
