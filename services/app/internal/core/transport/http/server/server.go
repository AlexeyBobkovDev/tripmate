package core_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"

	"github.com/AlexeyBobkovDev/tripmate/services/app/docs"
	core_logger "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/logger"
	core_middleware "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/transport/http/middleware"
)

type HTTPServer struct {
	mux         *http.ServeMux
	config      Config
	logger      *core_logger.Logger
	Middlewares []core_middleware.Middleware
}

func (s *HTTPServer) Handler() http.Handler {
	return core_middleware.ChainMiddleware(
		s.mux,
		s.Middlewares...,
	)
}

func (s *HTTPServer) Run(ctx context.Context) error {
	server := &http.Server{
		Addr:    s.config.Addr,
		Handler: s.Handler(),
	}

	errCh := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.logger.Debug("start HTTP server", zap.String("addr", s.config.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Debug("failed to start HTTP server", zap.Error(err))
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		shutDownCtx, shutDownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutDownCancel()

		if errShutdown := server.Shutdown(shutDownCtx); errShutdown != nil {
			errClose := server.Close()

			wg.Wait()

			if errClose != nil {
				combinedError := errors.Join(errShutdown, errClose)
				s.logger.Error("force shutdown failed", zap.Error(combinedError))
				return fmt.Errorf("force shutdown failed: %w", combinedError)
			}

			s.logger.Error("graceful shutdown failed", zap.Error(errShutdown))
			return fmt.Errorf("graceful shutdown failed: %w", errShutdown)
		}
		wg.Wait()
		s.logger.Debug("server stopped correctly")
	case err := <-errCh:
		wg.Wait()
		s.logger.Error("HTTP server stopped unexpectedly", zap.Error(err))
		return fmt.Errorf("HTTP server stopped unexpectedly: %w", err)
	}

	return nil
}

func (s *HTTPServer) RegisterRouters(routers ...*APIRouter) {
	for _, router := range routers {
		if router == nil || router.mux == nil {
			panic("cannot register nil router or mux")
		}

		path := "/api/" + router.APIVersion.ToString()

		routerMux := core_middleware.ChainMiddleware(router.mux, router.Middlewares...)

		s.mux.Handle(
			path+"/",
			http.StripPrefix(path, routerMux),
		)
	}
}

func (s *HTTPServer) RegisterHealth() {
	s.mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func (s *HTTPServer) RegisterSwagger() {
	s.mux.Handle(
		"/swagger/",
		httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"),
			httpSwagger.DefaultModelsExpandDepth(-1),
		),
	)

	s.mux.HandleFunc(
		"/swagger/doc.json",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(docs.SwaggerInfo.ReadDoc()))
		},
	)
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
