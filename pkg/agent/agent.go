package agent

import (
	"context"
	"log"

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
	collectors := collector.NewRunner(
		collector.NewPlayerInfo(url, client, prometheus.DefaultRegisterer, l),
	)
	collectors.Run(context.Background())
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
	return <-errCh // pointless channel
}
