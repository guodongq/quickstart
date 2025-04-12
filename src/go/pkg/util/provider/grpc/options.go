package grpc

import (
	"context"
	"fmt"
	"runtime/debug"

	logger "github.com/guodongq/quickstart/pkg/util/log"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type GrpcOptions struct {
	Port                        int
	HealthEnabled               bool
	RecoveryHandlerFunc         recovery.RecoveryHandlerFunc
	InterceptorLogger           logging.LoggerFunc
	AuthFunc                    auth.AuthFunc
	IgnoreProtoValidateMessages []protoreflect.MessageType
	UnaryInterceptors           []grpc.UnaryServerInterceptor
	StreamInterceptors          []grpc.StreamServerInterceptor
}

func getDefaultGrpcOptions() GrpcOptions {
	return GrpcOptions{
		Port:                        3000,
		HealthEnabled:               false,
		RecoveryHandlerFunc:         nil,
		InterceptorLogger:           nil,
		AuthFunc:                    nil,
		IgnoreProtoValidateMessages: nil,
		UnaryInterceptors:           nil,
		StreamInterceptors:          nil,
	}
}

func WithGrpcOptionsPort(port int) func(*GrpcOptions) {
	return func(options *GrpcOptions) {
		options.Port = port
	}
}

func WithGrpcOptionsHealthEnabled() func(*GrpcOptions) {
	return func(options *GrpcOptions) {
		options.HealthEnabled = true
	}
}

func WithGrpcOptionsRecoveryHandlerFunc(f recovery.RecoveryHandlerFunc) func(*GrpcOptions) {
	return func(options *GrpcOptions) {
		options.RecoveryHandlerFunc = f
	}
}

func WithGrpcOptionsDefaultRecoveryHandlerFunc() func(*GrpcOptions) {
	return func(options *GrpcOptions) {
		options.RecoveryHandlerFunc = func(p any) (err error) {
			logger.
				WithError(fmt.Errorf("%v", p)).
				WithField("stack", string(debug.Stack())).
				Error("panic recovered")

			return status.Errorf(codes.Internal, "%v", p)
		}
	}
}

func WithGrpcOptionsInterceptorLogger(f logging.LoggerFunc) func(*GrpcOptions) {
	return func(options *GrpcOptions) {
		options.InterceptorLogger = f
	}
}

func WithGrpcOptionsInterceptorLogging() func(*GrpcOptions) {
	return func(options *GrpcOptions) {
		options.InterceptorLogger = func(_ context.Context, level logging.Level, msg string, fields ...any) {
			f := make(logger.Fields)
			for i := 0; i < len(fields); i += 2 {
				key := fields[i]
				value := fields[i+1]
				f[fmt.Sprintf("%v", key)] = value
			}

			logEntry := logger.DefaultLogger().WithFields(f)
			switch level {
			case logging.LevelDebug:
				logEntry.Debug(msg)
			case logging.LevelInfo:
				logEntry.Info(msg)
			case logging.LevelWarn:
				logEntry.Warn(msg)
			case logging.LevelError:
				logEntry.Error(msg)
			default:
				logEntry.Panicf("unknown level %v", level)
			}
		}
	}
}

func WithGrpcOptionsAuthFunc(f auth.AuthFunc) func(*GrpcOptions) {
	return func(options *GrpcOptions) {
		options.AuthFunc = f
	}
}

func WithGrpcOptionsIgnoreProtoValidateMessages(msgs ...protoreflect.MessageType) func(*GrpcOptions) {
	return func(options *GrpcOptions) {
		options.IgnoreProtoValidateMessages = msgs
	}
}

func WithGrpcOptionsUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) func(*GrpcOptions) {
	return func(options *GrpcOptions) {
		options.UnaryInterceptors = append(options.UnaryInterceptors, interceptors...)
	}
}

func WithGrpcOptionsStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) func(*GrpcOptions) {
	return func(options *GrpcOptions) {
		options.StreamInterceptors = append(options.StreamInterceptors, interceptors...)
	}
}
