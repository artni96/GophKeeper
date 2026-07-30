package interceptors

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func PanicInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if recovered := recover(); recovered != nil {
				logger.Info("got panic",
					zap.String("error message", fmt.Sprintf("panic recovered: %v\n", recovered)),
					zap.String("call stack", string(debug.Stack())),
				)
				resp = nil
				err = status.Errorf(codes.Internal, "Internal Server Error")
			}
		}()
		return handler(ctx, req)
	}
}

func RequestLoggerInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		start := time.Now()
		resp, err = handler(ctx, req)
		respStatusCode := codes.OK.String()
		if err != nil {
			respStatusCode = status.Code(err).String()
		}

		duration := time.Since(start)

		logger.Info("Request done",
			zap.String("method", info.FullMethod),
			zap.Duration("duration", duration),
			zap.String("response_status", respStatusCode),
		)
		return resp, err
	}
}
