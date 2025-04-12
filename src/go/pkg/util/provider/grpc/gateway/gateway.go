package gateway

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	logger "github.com/guodongq/quickstart/pkg/util/log"
	"github.com/guodongq/quickstart/pkg/util/provider"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/guodongq/quickstart/pkg/util/provider/app"
	grpcProvider "github.com/guodongq/quickstart/pkg/util/provider/grpc"
	"google.golang.org/grpc"
)

type Gateway struct {
	provider.AbstractRunProvider

	options      GatewayOptions
	appProvider  *app.App
	grpcProvider *grpcProvider.Server

	client *grpc.ClientConn
	srv    *http.Server
	mux    *runtime.ServeMux
}

func New(
	appProvider *app.App,
	grpcProvider *grpcProvider.Server,
	optionFuncs ...func(*GatewayOptions),
) *Gateway {
	defaultOptions := getDefaultGatewayOptions()
	options := &defaultOptions

	for _, optionFunc := range optionFuncs {
		optionFunc(options)
	}

	return &Gateway{
		appProvider:  appProvider,
		grpcProvider: grpcProvider,
		options:      *options,
	}
}

func (p *Gateway) Run() error {
	if !p.options.Enabled {
		logger.Infof("GRPC Gateway Provider is disabled")
		return nil
	}

	if err := provider.WaitForRunningProvider(p.grpcProvider, 2); err != nil {
		return err
	}

	basePath := p.appProvider.ParsePath()
	serverAddr := p.grpcProvider.Listener().Addr().String()
	addr := fmt.Sprintf(":%d", p.options.Port)

	logEntry := logger.WithFields(logger.Fields{
		"basePath":   basePath,
		"serverAddr": serverAddr,
		"addr":       addr,
	})

	conn, err := grpc.NewClient(
		serverAddr,
		// todo: add interceptors here
		// grpc.WithChainUnaryInterceptor(),
		// grpc.WithChainStreamInterceptor(),
	)
	if err != nil {
		logEntry.WithError(err).Errorf("GRPC Gateway could not connect to GRPC server")
		return err
	}

	p.client = conn
	p.mux = runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(p.options.IncomingHeaderMatcher),
	)

	p.srv = &http.Server{
		Addr:    addr,
		Handler: NewMuxWrapper(p.mux, basePath),
	}

	p.SetRunning(true)

	logEntry.Info("GRPC Gateway Provider launched")
	if err = p.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logEntry.WithError(err).Error("GRPC Gateway Provider launch failed")
		return err
	}

	return nil
}

func (p *Gateway) RegisterServices(
	functions ...func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error,
) error {
	if !p.options.Enabled {
		return nil
	}
	if err := provider.WaitForRunningProvider(p.grpcProvider, 2); err != nil {
		return err
	}

	for _, function := range functions {
		if err := function(context.Background(), p.mux, p.client); err != nil {
			return err
		}
	}
	return nil
}

func (p *Gateway) Close() error {
	if !p.options.Enabled || p.client == nil {
		return p.AbstractRunProvider.Close()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := p.srv.Shutdown(ctx); err != nil {
		logger.WithError(err).Error("Error while closing GRPC Gateway REST server")
		return err
	}
	if err := p.client.Close(); err != nil {
		logger.WithError(err).Error("Error while closing GRPC Gateway connection to server")
		return err
	}

	return p.AbstractRunProvider.Close()
}

type MuxWrapper struct {
	mux      *runtime.ServeMux
	basePath string
}

func NewMuxWrapper(mux *runtime.ServeMux, basePath string) *MuxWrapper {
	return &MuxWrapper{
		mux:      mux,
		basePath: basePath,
	}
}

func (m *MuxWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = "/" + strings.TrimPrefix(r.URL.Path, m.basePath)
	m.mux.ServeHTTP(w, r)
}
