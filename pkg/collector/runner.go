package collector

import "context"

type Collector interface {
	scrape(context.Context) error
}

type Runner struct {
	cs []Collector
}

func NewRunner(cs ...Collector) *Runner {
	return &Runner{
		cs: cs,
	}
}

func (r *Runner) Run(ctx context.Context) error {
	// TODO
	for _, c := range r.cs {
		if err := c.scrape(ctx); err != nil {
			return err
		}
	}
	return nil
}
