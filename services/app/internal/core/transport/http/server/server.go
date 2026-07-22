package core_server

import (
	"context"
	"net/http"

	core_logger "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/logger"
	core_middleware "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/transport/http/middleware"
	"go.uber.org/zap"
)

type HTTPServer struct {
	mux         *http.ServeMux
	config      Config
	logger      *core_logger.Logger
	Middlewares []core_middleware.Middleware
}

func (s *HTTPServer) Run(ctx context.Context) {
	mux := core_middleware.ChainMiddleware(s.mux, s.Middlewares...)

	server := http.Server{
		Addr:    s.config.Addr,
		Handler: mux,
	}

	errCh := make(chan error, 1)
	go func() {
		s.logger.Debug("start HTTP server", zap.String("addr", s.config.Addr))
		if err := server.ListenAndServe(); err != nil {
			s.logger.Debug("failed to start HTTP server", zap.Error(err))
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		if err := server.Shutdown(ctx); err != nil {
			s.logger.Fatal("failed to shutdown server", zap.Error(err))
		}
		s.logger.Debug("server stopped correctly")
	case err := <-errCh:
		s.logger.Fatal("failed to start server", zap.Error(err))
	}
}

func (s *HTTPServer) RegisterRouters(routers ...*APIRouter) {
	for _, router := range routers {
		path := "/api/" + router.APIVersion.ToString()
		routerMux := core_middleware.ChainMiddleware(router.mux, router.Middlewares...)
		s.mux.Handle(
			path+"/",
			http.StripPrefix(path, routerMux),
		)
	}
}

func (s *HTTPServer) Health() {
	s.mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func NewHTTPServer(
	config Config,
	logger *core_logger.Logger,
	middlewares ...core_middleware.Middleware,
) *HTTPServer {
	return &HTTPServer{
		mux:         http.NewServeMux(),
		config:      config,
		logger:      logger,
		Middlewares: middlewares,
	}
}
