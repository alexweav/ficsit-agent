package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type PlayerInfo struct {
	baseURL *url.URL
	client  http.Client
	metrics *playerMetrics
	log     *log.Logger
}

func NewPlayerInfo(url *url.URL, client http.Client, reg prometheus.Registerer, log *log.Logger) *PlayerInfo {
	return &PlayerInfo{
		baseURL: url,
		client:  client,
		metrics: newPlayerMetrics(reg),
		log:     log,
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

func (p *PlayerInfo) scrape(ctx context.Context) error {
	uri := p.baseURL.JoinPath("/getPlayer")
	req, err := http.NewRequest(http.MethodGet, uri.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req = req.WithContext(ctx)

	p.log.Printf("executing request to %s", uri.String())
	resp, err := p.client.Do(req)
	if resp != nil {
		defer func() {
			if err := resp.Body.Close(); err != nil {
				p.log.Fatalf("failed to close response body: %s", err.Error())
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

	p.log.Printf("%v", players)

	return nil
}
