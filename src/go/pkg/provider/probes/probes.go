package probes

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/guodongq/quickstart/pkg/provider/app"

	logger "github.com/guodongq/quickstart/pkg/log"
	"github.com/guodongq/quickstart/pkg/provider"
)

type ProbeFunc func() error

type Probes struct {
	provider.AbstractRunProvider

	options         ProbesOptions
	appProvider     *app.App
	livenessProbes  []ProbeFunc
	readinessProbes []ProbeFunc

	srv *http.Server
}

func New(appProvider *app.App, optionsFuncs ...func(*ProbesOptions)) *Probes {
	defaultOptions := getDefaultProbesOptions()
	options := &defaultOptions

	for _, optionsFunc := range optionsFuncs {
		optionsFunc(options)
	}

	return &Probes{
		appProvider: appProvider,
		options:     *options,
	}
}

func (p *Probes) Init() error {
	if p.appProvider == nil {
		return fmt.Errorf("app Provider is required")
	}

	return nil
}

func (p *Probes) Run() error {
	if !p.options.Enabled {
		logger.Info("Probes Provider is disabled")
		return nil
	}

	addr := fmt.Sprintf(":%d", p.options.Port)
	livenessEndpoint := p.appProvider.ParseEndpoint(p.options.LivenessEndpoint)
	readinessEndpoint := p.appProvider.ParseEndpoint(p.options.ReadinessEndpoint)

	logEntry := logger.WithFields(logger.Fields{
		"addr":               addr,
		"liveness_endpoint":  livenessEndpoint,
		"readiness_endpoint": readinessEndpoint,
	})

	mux := http.NewServeMux()
	mux.HandleFunc(livenessEndpoint, p.livenessHandler)
	mux.HandleFunc(readinessEndpoint, p.readinessHandler)

	p.srv = &http.Server{Addr: addr, Handler: mux}
	p.SetRunning(true)

	logEntry.Info("Probes Provider Launched")
	if err := p.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logEntry.WithError(err).Error("Probes Provider launch failed")
		return err
	}

	return nil
}

func (p *Probes) Close() error {
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

func (p *Probes) livenessHandler(res http.ResponseWriter, req *http.Request) {
	reqDump, _ := httputil.DumpRequest(req, false)
	logger.WithField("req", string(reqDump)).Debug("Handling liveness request")
	for _, probe := range p.livenessProbes {
		if err := probe(); err != nil {
			res.WriteHeader(http.StatusServiceUnavailable)
			if _, err := res.Write([]byte(err.Error())); err != nil {
				logger.WithError(err).Warnf("Error while writing liveness data")
			}
			return
		}
	}
	res.WriteHeader(http.StatusOK)
}

func (p *Probes) readinessHandler(res http.ResponseWriter, req *http.Request) {
	reqDump, _ := httputil.DumpRequest(req, false)
	logger.WithField("req", string(reqDump)).Debug("Handling readiness request")
	for _, probe := range p.readinessProbes {
		if err := probe(); err != nil {
			res.WriteHeader(http.StatusServiceUnavailable)
			if _, err := res.Write([]byte(err.Error())); err != nil {
				logger.WithError(err).Warnf("Error while writing readiness data")
			}
			return
		}
	}
	res.WriteHeader(http.StatusOK)
}

func (p *Probes) AddLivenessProbes(fn ProbeFunc) {
	p.livenessProbes = append(p.livenessProbes, fn)
}

func (p *Probes) AddReadinessProbes(fn ProbeFunc) {
	p.readinessProbes = append(p.readinessProbes, fn)
}
