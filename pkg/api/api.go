package api

import (
	"context"
	"net/http"
	"time"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type API struct {
	log log.Logger
	srv http.Server
}

func New(l log.Logger) *API {
	return &API{
		log: l,
		srv: http.Server{
			Addr:         "127.0.0.1:1234",
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
	}
}

func (a *API) Run(ctx context.Context) error {
	a.register()
	return a.srv.ListenAndServe()
}

func (a *API) register() {
	http.HandleFunc("/", a.rootHandler)
	http.Handle("/metrics", promhttp.Handler())
}

func (a *API) rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("lol, lmao"))
}
