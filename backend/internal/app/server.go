package app

import (
	"context"
	"net/http"
	"time"

	infrahttp "go-shop-app-backend/internal/infra/http"
	"go-shop-app-backend/pkg/logger"
)

type Server struct {
	httpServer *http.Server
}

func NewHTTPServer(c *Container) *Server {
	router := infrahttp.NewRouter(c.DB, c.Config)

	srv := &http.Server{
		Addr:         ":" + c.Config.ServerPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{httpServer: srv}
}

func (s *Server) Start() error {
	logger.Info("server started", "addr", s.httpServer.Addr)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	logger.Info("server shutting down")
	return s.httpServer.Shutdown(ctx)
}
