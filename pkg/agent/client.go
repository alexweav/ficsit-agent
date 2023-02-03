package agent

import (
	"net/http"
	"net/url"
	"time"
)

func baseURL() *url.URL {
	url, _ := url.Parse("http://localhost:8080")
	return url
}

func newDefaultClient() http.Client {
	return http.Client{
		Timeout: 30 * time.Second,
	}
}
