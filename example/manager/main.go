package main

import (
	"context"

	"github.com/chmking/horde/manager"
)

func main() {
	ctx, _ := context.WithCancel(context.Background())

	m := &manager.Manager{}
	m.ListenAndServePublic()
	m.ListenAndServePrivate()
	m.Healthcheck(ctx)

	<-ctx.Done()
}
