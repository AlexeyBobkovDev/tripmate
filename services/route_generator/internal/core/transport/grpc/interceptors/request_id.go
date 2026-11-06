package core_interceptor

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func RequestIDInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		requestID := uuid.NewString()
		md := metadata.Pairs(
			"X-Request-ID", requestID,
		)
		ctx = metadata.NewIncomingContext(ctx, md)
		return handler(ctx, req)
	}
}
