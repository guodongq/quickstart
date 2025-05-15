package http

import (
	"context"
	"errors"
	"fmt"
	logger "github.com/guodongq/quickstart/pkg/log"
	"github.com/guodongq/quickstart/pkg/provider"
	"net/http"
	"time"
)

type HttpServer struct {
	provider.AbstractRunProvider

	srv *http.Server

	options HttpServerOptions
}

func New(optionsFuncs ...func(options *HttpServerOptions)) *HttpServer {
	defaultOptions := getDefaultHttpServerOptions()
	options := &defaultOptions

	for _, optionsFunc := range optionsFuncs {
		optionsFunc(options)
	}

	return &HttpServer{
		options: *options,
	}
}

func (h *HttpServer) Run() error {
	addr := fmt.Sprintf(":%d", h.options.Port)

	logEntry := logger.WithFields(logger.Fields{
		"addr": addr,
	})

	mux := http.NewServeMux()

	var handler http.Handler = mux
	if h.options.RegisterHandlers != nil {
		handler = h.options.RegisterHandlers(mux)
	}

	h.srv = &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	h.SetRunning(true)

	logEntry.Info("HTTP Server Provider Launched")
	if err := h.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logEntry.WithError(err).Error("HTTP Server Provider launch failed")
		return err
	}

	return nil
}

func (h *HttpServer) Close() error {
	if h.srv == nil {
		return h.AbstractRunProvider.Close()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	if err := h.srv.Shutdown(ctx); err != nil {
		logger.WithError(err).Error("Error while closing Probes server")
	}

	return h.AbstractRunProvider.Close()
}
