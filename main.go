package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("starting app")

	var handler http.HandlerFunc = func(rw http.ResponseWriter, req *http.Request) {
		log.Println("received request:", req)

		rw.Write([]byte("200 OK"))
	}

	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalln("starting http server error:", err.Error())
	}
}
