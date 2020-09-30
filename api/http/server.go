package http

import (
	"database/sql"
	"log"

	"github.com/tonytcb/bank-transactions-go/infra/repository"
	"github.com/tonytcb/bank-transactions-go/usecase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tonytcb/bank-transactions-go/api/http/handler"
)

// Server exposes the app through the HTTP protocol
type Server struct {
	logger  *log.Logger
	storage *sql.DB
}

// NewServer creates a Server struct with its dependencies
func NewServer(logger *log.Logger, storage *sql.DB) *Server {
	return &Server{logger: logger, storage: storage}
}

// Listen exposes the HTTP server running in the port 8080
func (s Server) Listen() {
	s.logger.Println("starting http server")

	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Output: s.logger.Writer()}))
	e.Use(middleware.Recover())

	e.POST("/accounts", s.createAccountHandler())

	s.logger.Fatalln(e.Start(":8080"))
}

func (s Server) createAccountHandler() echo.HandlerFunc {
	// @todo translate the standard response to use echo responder

	accountRepo := repository.NewAccount(s.storage)

	return func(ctx echo.Context) error {
		createAccount := handler.NewCreateAccount(
			s.logger,
			usecase.NewCreateAccount(accountRepo),
		)
		createAccount.Handler(ctx.Response().Writer, ctx.Request())

		return nil
	}
}
