package core_interceptor

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	core_logger "github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/logger"
)

func TraceInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		logger := core_logger.FromContext(ctx)
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Internal, "failed to retrieve metadata from context")
		}
		requestIDs := md.Get("X-Request-ID")
		if len(requestIDs) == 0 {
			return nil, status.Error(codes.Internal, "failed to retrieve X-Request-ID from metadata")
		}
		if len(requestIDs) > 1 {
			return nil, status.Error(codes.Internal, "multiple X-Request-ID keys")
		}

		requestID := requestIDs[0]

		logger.Debug(
			"start processing request",
			zap.String("requestID", requestID),
		)

		response, err := handler(ctx, req)
		if err != nil {
			logger.Debug(
				"stop processing request",
				zap.String("requestID", requestID),
				zap.Error(err),
			)
		} else {
			logger.Debug(
				"stop processing request",
				zap.String("requestID", requestID),
			)
		}
		return response, err
	}
}
