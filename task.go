package horde

import "context"

type TaskFunc func(ctx context.Context)

func (t TaskFunc) Exec(ctx context.Context) {
	t(ctx)
}

type Task struct {
	Name   string
	Func   TaskFunc
	Weight int
}
