package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/LegationPro/corpus-reader/internal/logger"
)

// Initialize the logger and set it globally
func initializeLogger() *slog.Logger {
	newLogger := logger.New()
	slog.SetDefault(newLogger)
	return newLogger
}

// Configuration for the server.
type Config struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type Server struct {
	server  *http.Server
	handler *http.ServeMux
	logger  *slog.Logger
	config  Config
}

func New(config Config) *Server {
	return &Server{
		server:  nil,
		handler: http.NewServeMux(),
		logger:  initializeLogger(),
		config:  config,
	}
}

func (s *Server) Start() {
	s.logger.Info("Starting server on address: " + s.config.Addr)

	s.server = &http.Server{
		Addr:         s.config.Addr,
		Handler:      http.NewServeMux(),
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
		IdleTimeout:  s.config.IdleTimeout,
	}

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Error(fmt.Sprintf("Failed to start server: %v", err))
	}
}

func (s *Server) Stop() {
	s.logger.Info("Shutting down server")

	if err := s.server.Shutdown(context.Background()); err != nil {
		s.logger.Error(fmt.Sprintf("Failed to shutdown server: %v", err))
	}
}
