package agent

import (
	"context"

	"github.com/alexweav/ficsit-agent/pkg/api"
	"github.com/alexweav/ficsit-agent/pkg/collector"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

type Agent struct {
	log        log.Logger
	api        *api.API
	collectors *collector.Runner
}

func New(cfg Config, l log.Logger) (*Agent, error) {
	client, err := collector.NewFRMClient(cfg.ModURL, l)
	if err != nil {
		return nil, err
	}
	collOps := collector.RunnerOpts{
		ScrapeInterval: cfg.ScrapeInterval,
		Log:            l,
	}
	collectors := collector.NewRunner(
		collOps,
		collector.NewForPlayers(client, prometheus.DefaultRegisterer, l),
		collector.NewForPower(client, prometheus.DefaultRegisterer, l),
	)
	api := api.New(l)
	return &Agent{
		log:        l,
		api:        api,
		collectors: collectors,
	}, nil
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
