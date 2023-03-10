package collector

import (
	"context"
	"fmt"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type PlayerCollector struct {
	client  *FRMClient
	metrics *playerMetrics
	log     log.Logger
}

func NewForPlayers(client *FRMClient, reg prometheus.Registerer, logger log.Logger) *PlayerCollector {
	return &PlayerCollector{
		client:  client,
		metrics: newPlayerMetrics(reg),
		log:     logger,
	}
}

type playerMetrics struct {
	Health *prometheus.GaugeVec
	Ping   *prometheus.GaugeVec
}

func newPlayerMetrics(r prometheus.Registerer) *playerMetrics {
	return &playerMetrics{
		Health: promauto.With(r).NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "satisfactory",
			Subsystem: "player",
			Name:      "health",
			Help:      "The current health of each player.",
		}, []string{"name"}),
		Ping: promauto.With(r).NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "satisfactory",
			Subsystem: "player",
			Name:      "ping",
			Help:      "The current ping time of the player in milliseconds.",
		}, []string{"name"}),
	}
}

type player struct {
	ID   string  `json:"ID"`
	Name string  `json:"PlayerName"`
	HP   float64 `json:"PlayerHP"`
	Ping int64   `json:"PingTime"`
}

func (p *PlayerCollector) scrape(ctx context.Context) error {
	players := make([]player, 0)
	err := p.client.GetJSON(ctx, "/getPlayer", &players)
	if err != nil {
		return fmt.Errorf("error fetching player data: %w", err)
	}

	p.metrics.Health.Reset()
	p.metrics.Ping.Reset()
	for _, pl := range players {
		if pl.Name != "" {
			p.metrics.Health.WithLabelValues(pl.Name).Set(pl.HP)
			p.metrics.Ping.WithLabelValues(pl.Name).Set(float64(pl.Ping))
		}
	}

	return nil
}
