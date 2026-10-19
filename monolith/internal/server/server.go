package server

import (
	"context"
	"log"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/spearszx/go-wallet/monolith/internal/middleware"
	"github.com/spearszx/go-wallet/monolith/internal/user/handlers"
	userService "github.com/spearszx/go-wallet/monolith/internal/user/service"
)

type Options struct {
	ListenAddr string
}

type Server struct {
	Options

	validator *validator.Validate

	logger     *slog.Logger
	httpServer *http.Server
}

func NewServer(config Options, logger *slog.Logger, validate *validator.Validate, userService userService.UserService) *Server {
	handler := http.NewServeMux()

	handler.Handle("/users",
		middleware.Recovery(
			middleware.Logger(
				logger,
				middleware.RequestID(
					middleware.ErrorHandler(
						middleware.MethodPost(
							handlers.Register(validate, userService),
						),
					),
				),
			),
		),
	)

	httpServer := http.Server{
		Addr:    config.ListenAddr,
		Handler: handler,
	}

	return &Server{
		Options:    config,
		logger:     logger,
		httpServer: &httpServer,
	}
}

func (s *Server) ListenAndServe() error {
	log.Println("running the server on port", s.ListenAddr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shitdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
