package prometheus

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	logger "github.com/guodongq/quickstart/pkg/log"
	"github.com/guodongq/quickstart/pkg/provider"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Prometheus struct {
	provider.AbstractRunProvider

	options PrometheusOptions
	srv     *http.Server
}

func New(optionsFuncs ...func(*PrometheusOptions)) *Prometheus {
	defaultOptions := getDefaultPrometheusOptions()
	options := &defaultOptions
	options.MergeIn(optionsFuncs...)

	return &Prometheus{options: *options}
}

func (p *Prometheus) Run() error {
	if !p.options.Enabled {
		logger.Info("Prometheus Provider is disabled")
		return nil
	}

	addr := fmt.Sprintf(":%d", p.options.Port)
	logEntry := logger.WithFields(logger.Fields{
		"addr":     addr,
		"endpoint": p.options.Endpoint,
	})

	mux := http.NewServeMux()
	mux.Handle(p.options.Endpoint, promhttp.Handler())

	p.srv = &http.Server{Addr: addr, Handler: mux}
	p.SetRunning(true)

	logEntry.Info("Prometheus Provider Launched")
	if err := p.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logEntry.WithError(err).Error("Prometheus Provider launch failed")
		return err
	}

	return nil
}

func (p *Prometheus) Close() error {
	if !p.options.Enabled || p.srv == nil {
		return p.AbstractRunProvider.Close()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	if err := p.srv.Shutdown(ctx); err != nil {
		logger.WithError(err).Error("Error while closing Probes server")
	}

	return p.AbstractRunProvider.Close()
}
