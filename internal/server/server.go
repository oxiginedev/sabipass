package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/oxiginedev/sabipass/config"
)

type Server struct {
	s      *http.Server
	stopFn func()
}

func NewServer(cfg *config.Config, stopFn func()) *Server {
	return &Server{
		s: &http.Server{
			Addr:              fmt.Sprintf(":%d", cfg.HTTP.Port),
			ReadTimeout:       10 * time.Second,
			WriteTimeout:      10 * time.Second,
			ReadHeaderTimeout: 2 * time.Second,
		},
		stopFn: stopFn,
	}
}

func (s *Server) Listen() {
	go func() {
		err := s.s.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("could not start server", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	s.waitForShutdown()
}

func (s *Server) SetHandler(handler http.Handler) {
	s.s.Handler = handler
}

func (s *Server) waitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server gracefully...\n")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.s.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("server exited properly")

	if s.stopFn != nil {
		s.stopFn()
	}
}
