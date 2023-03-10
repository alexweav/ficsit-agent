package collector

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-kit/log"
)

type Collector interface {
	scrape(context.Context) error
}

type RunnerOpts struct {
	ScrapeInterval time.Duration
	Log            log.Logger
}

type Runner struct {
	opts   RunnerOpts
	cs     []Collector
	ticker *time.Ticker
	log    log.Logger
}

func NewRunner(opts RunnerOpts, cs ...Collector) *Runner {
	return &Runner{
		opts:   opts,
		cs:     cs,
		ticker: time.NewTicker(opts.ScrapeInterval),
		log:    log.With(opts.Log, "component", "collector"),
	}
}

// Run starts the runner and blocks. It cannot be called concurrently.
func (r *Runner) Run(ctx context.Context) error {
	doneCh := make(chan bool) // TODO: not used, but we need a way to signal shutdowns eventually
	errCh := make(chan error)

	go func() {
		var wg sync.WaitGroup

		for {
			select {
			case <-doneCh:
				close(errCh)
				return
			case t := <-r.ticker.C:
				// TODO
				r.log.Log("msg", "Running all scrapers...")
				for _, c := range r.cs {
					coll := c

					wg.Add(1)
					go func() {
						defer wg.Done()
						if err := coll.scrape(ctx); err != nil {
							// TODO, log and continue, no need for errCh, horrible hack
							errCh <- fmt.Errorf("failed to scrape a target: %w", err)
						}
					}()

				}

				wg.Wait()
				r.log.Log("msg", "All scrapers for tick finished", "tick", t)
			}
		}
	}()

	return <-errCh
}
