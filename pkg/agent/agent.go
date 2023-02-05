package agent

import (
	"context"
	"log"
	"time"

	"github.com/alexweav/ficsit-agent/pkg/api"
	"github.com/alexweav/ficsit-agent/pkg/collector"
	"github.com/prometheus/client_golang/prometheus"
)

type Agent struct {
	log        *log.Logger
	api        *api.API
	collectors *collector.Runner
}

func New(l *log.Logger) *Agent {
	url := baseURL()
	client := newDefaultClient()
	collOps := collector.RunnerOpts{
		ScrapeInterval: 10 * time.Second,
		Log:            l,
	}
	collectors := collector.NewRunner(
		collOps,
		collector.NewForPlayers(url, client, prometheus.DefaultRegisterer, l),
		collector.NewForPower(url, client, prometheus.DefaultRegisterer, l),
	)
	api := api.New(l)
	return &Agent{
		log:        l,
		api:        api,
		collectors: collectors,
	}
}

func (a *Agent) Run(ctx context.Context) error {
	errCh := make(chan error)
	go func() {
		errCh <- a.api.Run(ctx)
	}()
	go func() {
		errCh <- a.collectors.Run(ctx)
	}()
	return <-errCh
}
