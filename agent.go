package goteem

import "context"

type Task struct {
	Name   string
	Weight int
	Func   func(ctx context.Context)
}

type Behavior struct {
	Tasks []*Task
}

type Agent struct {
	Behavior *Behavior
}

func (a *Agent) Teem() {
}
