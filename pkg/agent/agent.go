package agent

import (
	"context"
	"net/url"
	"time"

	"github.com/alexweav/ficsit-agent/pkg/api"
	"github.com/alexweav/ficsit-agent/pkg/collector"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

type Config struct {
	ScrapeInterval time.Duration
	ModURL         *url.URL
}

type Agent struct {
	log        log.Logger
	api        *api.API
	collectors *collector.Runner
}

func New(cfg Config, l log.Logger) *Agent {
	client := collector.NewFRMClient(cfg.ModURL, l)
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
