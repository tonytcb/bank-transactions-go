package main

import (
	"log"
	"os"

	"github.com/tonytcb/bank-transactions-go/api"
	"github.com/tonytcb/bank-transactions-go/api/http"
)

func main() {
	var (
		logger                = log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)
		httpServer api.Server = http.NewServer(logger)
	)

	logger.Println("starting app")

	httpServer.Listen()
}
