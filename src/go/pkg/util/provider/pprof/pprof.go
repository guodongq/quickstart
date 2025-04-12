package pprof

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/pprof"
	"time"

	logger "github.com/guodongq/quickstart/pkg/util/log"
	"github.com/guodongq/quickstart/pkg/util/provider"
)

type PProf struct {
	provider.AbstractRunProvider

	options PProfOptions
	srv     *http.Server
}

func New(optionFuncs ...func(*PProfOptions)) *PProf {
	defaultOptions := getDefaultOptions()
	options := &defaultOptions
	for _, optionFunc := range optionFuncs {
		optionFunc(options)
	}

	return &PProf{
		options: *options,
	}
}

func (p *PProf) Run() error {
	if !p.options.Enabled {
		logger.Info("PProf Provider is disabled")
		return nil
	}

	addr := fmt.Sprintf(":%d", p.options.Port)
	logEntry := logger.WithFields(logger.Fields{
		"addr":     addr,
		"endpoint": p.options.Endpoint,
	})

	mux := http.NewServeMux()
	mux.HandleFunc(p.options.Endpoint+"/", pprof.Index)
	mux.HandleFunc(p.options.Endpoint+"/cmdline", pprof.Cmdline)
	mux.HandleFunc(p.options.Endpoint+"/profile", pprof.Profile)
	mux.HandleFunc(p.options.Endpoint+"/symbol", pprof.Symbol)
	mux.HandleFunc(p.options.Endpoint+"/trace", pprof.Trace)

	p.srv = &http.Server{Addr: addr, Handler: mux}
	p.SetRunning(true)

	logEntry.Info("PProf Provider Launched")
	if err := p.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logEntry.WithError(err).Error("PProf Provider launch failed")
		return err
	}

	return nil
}

func (p *PProf) Close() error {
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
