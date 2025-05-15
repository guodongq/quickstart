package proxy

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	logger "github.com/guodongq/quickstart/pkg/log"
	"github.com/guodongq/quickstart/pkg/provider"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Proxy struct {
	provider.AbstractRunProvider

	options      ProxyOptions
	srv          *http.Server
	reverseProxy *httputil.ReverseProxy
}

func New(optionsFuncs ...func(*ProxyOptions)) *Proxy {
	defaultOptions := getDefaultProxyOptions()
	options := &defaultOptions
	options.MergeIn(optionsFuncs...)

	return &Proxy{
		options: *options,
	}
}

func (p *Proxy) Init() error {
	targetURL, err := url.Parse(p.options.TargetURL)
	if err != nil {
		logger.
			WithField("target_url", p.options.TargetURL).
			WithError(err).
			Errorf("%s Proxy Provider initialization failed", p.options.Prefix)
		return err
	}
	p.reverseProxy = httputil.NewSingleHostReverseProxy(targetURL)
	p.reverseProxy.Transport = &loggingTransport{opts: p.options}
	logger.
		WithField("target_url", targetURL).
		Infof("%s Proxy initialized", cases.Title(language.Und).String(p.options.Prefix))
	return nil
}

func (p *Proxy) Run() error {
	if !p.options.Enabled {
		logger.Infof("%s Proxy is disabled", cases.Title(language.Und).String(p.options.Prefix))
		return nil
	}

	addr := fmt.Sprintf(":%d", p.options.Port)

	logEntry := logger.WithFields(logger.Fields{
		"addr":     addr,
		"endpoint": p.options.Endpoint,
	})

	mux := http.NewServeMux()
	mux.HandleFunc(p.options.Endpoint, func(res http.ResponseWriter, req *http.Request) {
		req.Host = req.URL.Host
		p.reverseProxy.ServeHTTP(res, req)
	})

	p.srv = &http.Server{Addr: addr, Handler: mux}
	p.SetRunning(true)

	logEntry.Infof("%s Proxy launched", cases.Title(language.Und).String(p.options.Prefix))
	if err := p.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logEntry.WithError(err).Errorf("%s Proxy launch failed", cases.Title(language.Und).String(p.options.Prefix))
		return err
	}

	return nil
}

func (p *Proxy) Close() error {
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

type loggingTransport struct {
	opts ProxyOptions
}

func (t *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Log the reqBytes before sending.
	reqBytes, err := httputil.DumpRequestOut(req, t.opts.Debug)
	if err != nil {
		return nil, err
	}
	logEntry := logger.WithField("request", string(reqBytes))
	logEntry.Debugf("Performing proxy request to %s", t.opts.Prefix)

	// Perform the actual reqBytes.
	res, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		logEntry.WithError(err).Warnf("Received error from %s", t.opts.Prefix)
		return res, err
	}

	// Log the resBytes.
	resBytes, err := httputil.DumpResponse(res, t.opts.Debug)
	if err != nil {
		return nil, err
	}
	logEntry.WithField("response", string(resBytes)).Infof("Finished request to %s", t.opts.Prefix)

	return res, err
}
