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

type PowerCollector struct {
	baseURL *url.URL
	client  http.Client
	metrics *powerMetrics
	log     *log.Logger
}

func NewForPower(url *url.URL, client http.Client, reg prometheus.Registerer, log *log.Logger) *PowerCollector {
	return &PowerCollector{
		baseURL: url,
		client:  client,
		metrics: newPowerMetrics(reg),
		log:     log,
	}
}

type powerMetrics struct {
	Capacity       *prometheus.GaugeVec
	Production     *prometheus.GaugeVec
	MaxConsumption *prometheus.GaugeVec
	Consumption    *prometheus.GaugeVec
	Tripped        *prometheus.GaugeVec
}

func newPowerMetrics(r prometheus.Registerer) *powerMetrics {
	return &powerMetrics{
		Capacity: promauto.With(r).NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "satisfactory",
			Subsystem: "power",
			Name:      "capacity",
			Help:      "The theoretical maximum amount of power that can be produced by all machines in the network.",
		}, []string{"network"}),
		Production: promauto.With(r).NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "satisfactory",
			Subsystem: "power",
			Name:      "production",
			Help:      "The current amount of power being produced.",
		}, []string{"network"}),
		MaxConsumption: promauto.With(r).NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "satisfactory",
			Subsystem: "power",
			Name:      "max_consumption",
			Help:      "The theoretical maximum amount of power that can be consumed by all machines in the network.",
		}, []string{"network"}),
		Consumption: promauto.With(r).NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "satisfactory",
			Subsystem: "power",
			Name:      "consumption",
			Help:      "The current amount of power being consumed.",
		}, []string{"network"}),
		Tripped: promauto.With(r).NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "satisfactory",
			Subsystem: "power",
			Name:      "tripped",
			Help:      "Whether the current network has a tripped breaker.",
		}, []string{"network"}),
	}
}

func (p *powerMetrics) reset() {
	p.Capacity.Reset()
	p.Production.Reset()
	p.MaxConsumption.Reset()
	p.Consumption.Reset()
	p.Tripped.Reset()
}

func (p *powerMetrics) record(n network) {
	id := fmt.Sprint(n.CircuitID)
	p.Capacity.WithLabelValues(id).Set(n.PowerCapacity)
	p.Production.WithLabelValues(id).Set(n.PowerProduction)
	p.MaxConsumption.WithLabelValues(id).Set(n.PowerMaxConsumed)
	p.Consumption.WithLabelValues(id).Set(n.PowerConsumed)
	tripped := 0.0
	if n.FuseTriggered {
		tripped = 1.0
	}
	p.Tripped.WithLabelValues(id).Set(tripped)
}

type network struct {
	CircuitID        int     `json:"CircuitID"`
	PowerCapacity    float64 `json:"PowerCapacity"`
	PowerProduction  float64 `json:"PowerProduction"`
	PowerConsumed    float64 `json:"PowerConsumed"`
	PowerMaxConsumed float64 `json:"PowerMaxConsumed"`
	FuseTriggered    bool    `json:"FuseTriggered"`
}

func (p *PowerCollector) scrape(ctx context.Context) error {
	uri := p.baseURL.JoinPath("/getPower")
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
	networks := make([]network, 0)
	err = json.Unmarshal(data, &networks)
	if err != nil {
		return fmt.Errorf("error deserializing response: %w", err)
	}

	p.metrics.reset()
	for _, n := range networks {
		p.metrics.record(n)
	}

	return nil
}
