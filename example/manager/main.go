package main

import (
	"context"
	"log"

	"github.com/chmking/horde/manager"
)

func main() {
	ctx, _ := context.WithCancel(context.Background())

	m := manager.New()
	log.Fatal(m.Listen(ctx))
}
