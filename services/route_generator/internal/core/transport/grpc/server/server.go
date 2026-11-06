package core_server

import (
	"context"
	"fmt"
	"net"
	"sync"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	core_logger "github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/logger"
)

type GRPCServer struct {
	*grpc.Server
	config       Config
	logger       *core_logger.Logger
	Interceptors *[]grpc.UnaryServerInterceptor
}

func (s *GRPCServer) Run(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.config.Addr)
	if err != nil {
		return fmt.Errorf("failed to listen tcp port %s: %w", s.config.Addr, err)
	}

	server := s.Server

	errCh := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.logger.Debug("start gRPC server", zap.String("addr", s.config.Addr))
		if err := server.Serve(listener); err != nil {
			s.logger.Debug("failed to start gRPC server", zap.Error(err))
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		server.GracefulStop()

		wg.Wait()

		s.logger.Debug("server stopped correctly")
	case err := <-errCh:
		wg.Wait()
		s.logger.Error("gRPC server stopped unexpectedly", zap.Error(err))
		return fmt.Errorf("gRPC server stopped unexpectedly: %w", err)
	}

	return nil
}

func (s *GRPCServer) RegisterHealth() {
	healthServer := health.NewServer()
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(s, healthServer)
}

func NewGRPCServer(
	config Config,
	logger *core_logger.Logger,
	interceptors ...grpc.UnaryServerInterceptor,
) *GRPCServer {
	return &GRPCServer{
		Server: grpc.NewServer(
			grpc.ChainUnaryInterceptor(interceptors...),
		),
		config:       config,
		logger:       logger,
		Interceptors: &interceptors,
	}
}
