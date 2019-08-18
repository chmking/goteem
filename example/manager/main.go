package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"

	"github.com/chmking/horde/manager"
)

func main() {
	go func() {
		log.Fatal(http.ListenAndServe(":6061", nil))
	}()

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		m := manager.New()
		log.Fatal(m.Listen(ctx))
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	s := <-c
	fmt.Println("Got signal:", s)
	cancel()
}
