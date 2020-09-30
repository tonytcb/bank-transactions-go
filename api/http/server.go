package http

import (
	"log"
	"math/rand"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tonytcb/bank-transactions-go/api/http/handler"
)

// Server exposes the app through the HTTP protocol
type Server struct {
	logger *log.Logger
}

// NewServer creates a Server struct with its dependencies
func NewServer(logger *log.Logger) *Server {
	return &Server{logger: logger}
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

	return func(ctx echo.Context) error {
		createAccount := handler.NewCreateAccount(
			s.logger,
			&fakeAccountCreator{},
		)
		createAccount.Handler(ctx.Response().Writer, ctx.Request())

		return nil
	}
}

type fakeAccountCreator struct {
}

func (a fakeAccountCreator) Create(_ string) (int, error) {
	rand.Seed(time.Now().Unix())

	return int(rand.Int31n(100)), nil
}
