package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type PlayerCollector struct {
	baseURL *url.URL
	client  http.Client
	metrics *playerMetrics
	log     log.Logger
}

func NewForPlayers(url *url.URL, client http.Client, reg prometheus.Registerer, logger log.Logger) *PlayerCollector {
	return &PlayerCollector{
		baseURL: url,
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
	ID   string `json:"ID"`
	Name string `json:"PlayerName"`
	HP   int    `json:"PlayerHP"`
	Ping int64  `json:"PingTime"`
}

func (p *PlayerCollector) scrape(ctx context.Context) error {
	uri := p.baseURL.JoinPath("/getPlayer")
	req, err := http.NewRequest(http.MethodGet, uri.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req = req.WithContext(ctx)

	p.log.Log("msg", "Executing request", "url", uri.String())
	resp, err := p.client.Do(req)
	if resp != nil {
		defer func() {
			if err := resp.Body.Close(); err != nil {
				p.log.Log("msg", "Failed to close response body", "err", err)
			}
		}()
	}
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}
	players := make([]player, 0)
	err = json.Unmarshal(data, &players)
	if err != nil {
		return fmt.Errorf("error deserializing response: %w", err)
	}

	p.metrics.Health.Reset()
	p.metrics.Ping.Reset()
	for _, pl := range players {
		if pl.Name != "" {
			p.metrics.Health.WithLabelValues(pl.Name).Set(float64(pl.HP))
			p.metrics.Ping.WithLabelValues(pl.Name).Set(float64(pl.Ping))
		}
	}

	return nil
}
