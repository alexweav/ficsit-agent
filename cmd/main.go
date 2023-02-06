package main

import (
	"context"
	"flag"
	"net/url"
	"os"
	"os/signal"

	"github.com/alexweav/ficsit-agent/pkg/agent"
	"github.com/go-kit/log"
)

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	flags := flag.NewFlagSet("ficsit-agent", flag.PanicOnError)
	cfg := agent.Config{}
	cfg.RegisterFlags(flags)
	flags.Parse(os.Args[1:])

	agent, err := agent.New(cfg, logger)
	if err != nil {
		logger.Log("msg", "Failed to initialize agent", "error", err)
	}

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

func baseURL() *url.URL {
	url, _ := url.Parse("http://localhost:8080")
	return url
}
