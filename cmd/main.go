package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/alexweav/ficsit-agent/pkg/agent"
	"github.com/go-kit/log"
)

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))

	agent := agent.New(logger)

	inter := make(chan os.Signal, 1)
	signal.Notify(inter, os.Interrupt)

	errCh := make(chan error)
	go func() {
		errCh <- agent.Run(context.Background())
	}()

	select {
	case err := <-errCh:
		if err != nil {
			logger.Log("msg", "Exited with fatal error", "error", err)
		}
	// Not great, doesn't wait on anything to close
	case <-inter:
		logger.Log("msg", "Exiting...")
		os.Exit(0)
	}

}
