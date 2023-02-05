package collector

import (
	"context"
	"fmt"
	"log"
	"time"
)

type Collector interface {
	scrape(context.Context) error
}

type RunnerOpts struct {
	ScrapeInterval time.Duration
	Log            *log.Logger
}

type Runner struct {
	opts   RunnerOpts
	cs     []Collector
	ticker *time.Ticker
	log    *log.Logger
}

func NewRunner(opts RunnerOpts, cs ...Collector) *Runner {
	return &Runner{
		opts:   opts,
		cs:     cs,
		ticker: time.NewTicker(opts.ScrapeInterval),
		log:    opts.Log,
	}
}

// Run starts the runner and blocks. It cannot be called concurrently.
func (r *Runner) Run(ctx context.Context) error {
	doneCh := make(chan bool) // TODO: not used, but we need a way to signal shutdowns eventually
	errCh := make(chan error)

	go func() {
		for {
			select {
			case <-doneCh:
				close(errCh)
				return
			case <-r.ticker.C:
				// TODO
				for _, c := range r.cs {
					// TODO do not serialize scrapes, use waitgroup
					r.log.Println("Scraping a target...")
					if err := c.scrape(ctx); err != nil {
						// TODO, log and continue, no need for errCh, horrible hack
						errCh <- fmt.Errorf("failed to scrape a target: %w", err)
						break
					}
				}
			}
		}
	}()

	return <-errCh
}
