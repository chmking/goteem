package main

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/chmking/horde/logger"
	"github.com/chmking/horde/logger/log"
	"github.com/chmking/horde/manager/service"
)

func main() {
	log.Logger = logger.NewZeroLogger(logger.NewZeroConsoleWriter(os.Stderr))

	go func() {
		err := http.ListenAndServe(":6061", nil)
		log.Fatal().Err(err).Msg("pprof quit unexpectedly")
	}()

	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		s := service.New()
		if err := s.Listen(ctx); err != nil {
			log.Error().Err(err).Msg("manager quit unexpectedly")
		}
		c <- syscall.SIGQUIT
	}()

	sig := <-c
	log.Info().Msg("Got signal: " + sig.String())
	cancel()
}
