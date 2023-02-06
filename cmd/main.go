package main

import (
	"context"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/alexweav/ficsit-agent/pkg/agent"
	"github.com/go-kit/log"
)

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	cfg := agent.Config{
		ScrapeInterval: 10 * time.Second,
		ModURL:         baseURL(),
	}
	agent := agent.New(cfg, logger)

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
