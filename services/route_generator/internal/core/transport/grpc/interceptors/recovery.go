package core_interceptor

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	core_logger "github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/logger"
)

func RecoveryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		defer func() {
			logger := core_logger.FromContext(ctx)

			if p := recover(); p != nil {
				logger.Error(
					"panic occurred",
					zap.Any("panic", p),
					zap.Stack("stack"),
				)
				resp = nil
				err = status.Error(codes.Internal, "internal server error")

			}
		}()

		return handler(ctx, req)
	}
}
