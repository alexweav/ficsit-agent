package agent

import (
	"context"
	"log"

	"github.com/alexweav/ficsit-agent/pkg/api"
)

type Agent struct {
	log *log.Logger
	api *api.API
}

func New(l *log.Logger) *Agent {
	return &Agent{
		log: l,
		api: api.New(l),
	}
}

func (a *Agent) Run(ctx context.Context) error {
	errCh := make(chan error)
	go func() {
		errCh <- a.api.Run(ctx)
	}()
	return <-errCh // pointless channel
}
