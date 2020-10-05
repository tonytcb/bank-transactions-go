package http

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tonytcb/bank-transactions-go/api/http/handler"
	stdmiddleware "github.com/tonytcb/bank-transactions-go/api/http/middleware"
	"github.com/tonytcb/bank-transactions-go/infra/repository"
	"github.com/tonytcb/bank-transactions-go/usecase"
)

// Server exposes the app through the HTTP protocol
type Server struct {
	logger  *log.Logger
	storage *sql.DB
	port    int
}

// NewServer creates a Server struct with its dependencies
func NewServer(logger *log.Logger, storage *sql.DB, port int) *Server {
	return &Server{logger: logger, storage: storage, port: port}
}

// Listen exposes the HTTP server running in the port 8080
func (s Server) Listen() {
	s.logger.Println("starting http server")

	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(s.middleware(stdmiddleware.NewRequestID().Handler))
	e.Use(s.middleware(stdmiddleware.NewLogger(s.logger).Handler))

	e.POST("/accounts", s.createAccountHandler())
	e.GET("/accounts/:id", s.findAccountByIDHandler())
	e.POST("/transactions", s.createTransactionHandler())

	s.logger.Fatalln(e.Start(fmt.Sprintf(":%d", s.port)))
}

func (s Server) createAccountHandler() echo.HandlerFunc {
	createAccount := handler.NewCreateAccount(
		s.logger,
		usecase.NewCreateAccount(repository.NewAccountWriter(s.storage)),
	)

	return s.handler(createAccount.Handler)
}

func (s Server) findAccountByIDHandler() echo.HandlerFunc {
	findAccount := handler.NewFindAccount(
		s.logger,
		usecase.NewFindAccount(repository.NewAccountReader(s.storage)),
	)

	return s.handler(findAccount.Handler)
}

func (s Server) createTransactionHandler() echo.HandlerFunc {
	createTransaction := handler.NewCreateTransaction(
		s.logger,
		usecase.NewCreateTransaction(repository.NewTransaction(s.storage)),
	)

	return s.handler(createTransaction.Handler)
}

// handler translates a standard http handler to an echo handler
func (s Server) handler(fn func(http.ResponseWriter, *http.Request)) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		fn(ctx.Response().Writer, ctx.Request())

		return nil
	}
}

// middleware translates a standard http middleware to an echo middleware
func (s Server) middleware(fn func(http.ResponseWriter, *http.Request, http.HandlerFunc)) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			var nextFn http.HandlerFunc = func(rw http.ResponseWriter, r *http.Request) {
				// todo fix translation
				next(ctx)
			}

			fn(ctx.Response().Writer, ctx.Request(), nextFn)

			return nil
		}
	}
}
