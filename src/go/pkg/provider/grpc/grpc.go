package grpc

import (
	"fmt"
	"github.com/bufbuild/protovalidate-go"
	"net"

	logger "github.com/guodongq/quickstart/pkg/log"
	"github.com/guodongq/quickstart/pkg/provider"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	provider.AbstractRunProvider

	options  GrpcOptions
	server   *grpc.Server
	listener net.Listener
}

func New(optionFuncs ...func(*GrpcOptions)) *Server {
	defaultOptions := getDefaultGrpcOptions()
	options := &defaultOptions

	for _, optionFunc := range []func(*GrpcOptions){
		WithGrpcOptionsDefaultRecoveryHandlerFunc(),
		WithGrpcOptionsInterceptorLogging(),
	} {
		optionFunc(options)
	}

	for _, optionFunc := range optionFuncs {
		optionFunc(options)
	}

	return &Server{
		options: *options,
	}
}

func (p *Server) Init() error {
	validator, err := protovalidate.New()
	if err != nil {
		return err
	}

	unaryInterceptors := []grpc.UnaryServerInterceptor{
		recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(p.options.RecoveryHandlerFunc)),
		logging.UnaryServerInterceptor(p.options.InterceptorLogger),
		auth.UnaryServerInterceptor(p.options.AuthFunc),
		protovalidate_middleware.UnaryServerInterceptor(
			validator,
			protovalidate_middleware.WithIgnoreMessages(p.options.IgnoreProtoValidateMessages...),
		),
	}
	streamInterceptors := []grpc.StreamServerInterceptor{
		recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(p.options.RecoveryHandlerFunc)),
		logging.StreamServerInterceptor(p.options.InterceptorLogger),
		auth.StreamServerInterceptor(p.options.AuthFunc),
		protovalidate_middleware.StreamServerInterceptor(
			validator,
			protovalidate_middleware.WithIgnoreMessages(p.options.IgnoreProtoValidateMessages...),
		),
	}

	p.server = grpc.NewServer(
		grpc.ChainUnaryInterceptor(append(unaryInterceptors, p.options.UnaryInterceptors...)...),
		grpc.ChainStreamInterceptor(append(streamInterceptors, p.options.StreamInterceptors...)...),
	)
	return nil
}

func (p *Server) Run() error {
	addr := fmt.Sprintf(":%d", p.options.Port)
	reflection.Register(p.server)

	logEntry := logger.WithField("addr", addr)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logEntry.WithError(err).Error("GRPC engine Listener could not be created")
		return err
	}
	p.listener = listener
	p.SetRunning(true)
	p.registerHealthEndpoint()

	logEntry.Info("GRPC engine Provider launched")
	if err = p.server.Serve(listener); err != nil {
		logEntry.WithError(err).Error("GRPC engine Provider launch failed")
		return err
	}

	return nil
}

func (p *Server) Listener() net.Listener {
	return p.listener
}

func (p *Server) Close() error {
	p.server.GracefulStop()

	return p.AbstractRunProvider.Close()
}

func (p *Server) registerHealthEndpoint() {
	if !p.options.HealthEnabled {
		logger.Debugf("GRPC engine Health Check is disabled")
		return
	}

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(p.server, healthServer)
	logger.Debugf("GRPC engine Health Check registered")
}
