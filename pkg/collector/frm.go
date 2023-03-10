package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/go-kit/log"
)

type FRMClient struct {
	url    *url.URL
	client http.Client
	logger log.Logger
}

func NewFRMClient(baseURL string, l log.Logger) (*FRMClient, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	return &FRMClient{
		url: u,
		client: http.Client{
			Timeout: 30 * time.Second,
		},
		logger: l,
	}, nil
}

func (c *FRMClient) GetJSON(ctx context.Context, uri string, into any) error {
	url := c.url.JoinPath(uri)
	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req = req.WithContext(ctx)

	c.logger.Log("msg", "Executing request", "url", url.String())
	resp, err := c.client.Do(req)
	if resp != nil {
		defer func() {
			if err := resp.Body.Close(); err != nil {
				c.logger.Log("msg", "Failed to close response body", "err", err)
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
	err = json.Unmarshal(data, into)
	if err != nil {
		return fmt.Errorf("error deserializing response: %w", err)
	}

	return nil
}
