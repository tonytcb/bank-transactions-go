package main

import (
	"database/sql"
	"errors"
	"log"
	"os"

	"github.com/tonytcb/bank-transactions-go/api"
	"github.com/tonytcb/bank-transactions-go/api/http"
	"github.com/tonytcb/bank-transactions-go/infra/storage"
)

func main() {
	// todo improve the logger struct with common methods (INFO, WARN, ERROR, ...) and a way to track logs through the same process

	var logger = log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)

	logger.Println("starting app")

	db, err := newStorage()
	if err != nil {
		logger.Fatalln("error to start storage:", err.Error())
		return
	}
	defer db.Close()

	var httpServer api.Server = http.NewServer(logger, db)

	httpServer.Listen()
}

func newStorage() (*sql.DB, error) {
	if os.Getenv("MYSQL_HOST") == "" {
		return nil, errors.New("storage environment credentials not defined")
	}

	return storage.NewMySQLConnection(storage.NewConfig(
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_DATABASE"),
		os.Getenv("MYSQL_USER"),
	))
}
