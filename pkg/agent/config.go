package agent

import (
	"flag"
	"time"
)

type Config struct {
	ScrapeInterval time.Duration
	ModURL         string
}

func (c *Config) RegisterFlags(fs *flag.FlagSet) {
	fs.DurationVar(&c.ScrapeInterval, "scrape-interval", 10*time.Second, "rate at which to scrape data")
	fs.StringVar(&c.ModURL, "frm-url", "http://localhost:8080", "the URL of the ficsit remote management mod")
}
