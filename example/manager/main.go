package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/chmking/horde/manager"
)

func main() {
	go func() {
		log.Fatal(http.ListenAndServe(":6061", nil))
	}()

	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		m := manager.New()
		if err := m.Listen(ctx); err != nil {
			log.Print(err)
		}
		c <- syscall.SIGQUIT
	}()

	s := <-c
	fmt.Println("Got signal:", s)
	cancel()
}
