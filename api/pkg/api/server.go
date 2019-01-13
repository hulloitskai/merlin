package api

import (
	"context"
	"net/http"

	"go.uber.org/zap"

	"github.com/stevenxie/merlin/api/pkg/api/routes"
	"github.com/stevenxie/merlin/api/pkg/scrape"
	ess "github.com/unixpickle/essentials"
)

// A Server serves a REST API for accessing data from EDGAR.
type Server struct {
	srv *http.Server
	l   *zap.SugaredLogger
}

// NewServer makes a new Server.
func NewServer(logger *zap.SugaredLogger) (*Server, error) {
	if logger == nil {
		logger = zap.NewNop().Sugar()
	}

	// Make and configure router.
	cfg := routes.Config{
		Scraper: scrape.NewScraper(),
		Logger:  logger.Named("routes"),
	}
	router, err := routes.NewRouter(&cfg)
	if err != nil {
		return nil, ess.AddCtx("api: creating router", err)
	}

	srv := &http.Server{Handler: router}
	return &Server{
		srv: srv,
		l:   logger,
	}, nil
}

// ListenAndServe starts the server on the specified address.
func (s *Server) ListenAndServe(addr string) error {
	s.srv.Addr = addr
	return s.srv.ListenAndServe()
}

// Shutdown gracefully shuts down the Server, closing all existing connections.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
