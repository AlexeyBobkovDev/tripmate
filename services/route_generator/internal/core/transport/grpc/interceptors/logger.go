package core_interceptor

import (
	"context"

	"google.golang.org/grpc"

	core_logger "github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/logger"
)

func LoggerInterceptor(logger *core_logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		ctx = core_logger.ToContext(ctx, logger)

		return handler(ctx, req)
	}
}
