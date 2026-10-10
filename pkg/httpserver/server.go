package httpserver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

type Server struct {
	eg errgroup.Group

	server *http.Server
	notify chan error

	shutdownTimeout time.Duration

	logger *slog.Logger
}

func New(l *slog.Logger, address string, handler http.Handler) *Server {

	s := &Server{
		server: &http.Server{
			Addr:              address,
			Handler:           handler,
			ReadHeaderTimeout: 5 * time.Second,
			IdleTimeout:       60 * time.Second,
		},
		shutdownTimeout: 5 * time.Second,
		notify:          make(chan error, 1),
		logger:          l,
	}
	return s
}

func (s *Server) Start() {
	s.eg.Go(func() error {
		defer close(s.notify)

		err := s.server.ListenAndServe()

		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}

		if err != nil {
			err = fmt.Errorf("httpserver - Start: %w", err)
		}

		s.notify <- err
		return err
	})
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		s.shutdownTimeout,
	)
	defer cancel()

	shutdownErr := s.server.Shutdown(ctx)

	var closeErr error
	if shutdownErr != nil {
		closeErr = s.server.Close()
	}
	serveErr := s.eg.Wait()

	return errors.Join(shutdownErr, closeErr, serveErr)
}
