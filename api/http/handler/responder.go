package handler

import "net/http"

type responder struct {
	rw http.ResponseWriter
}

func newResponder(rw http.ResponseWriter) *responder {
	return &responder{rw: rw}
}

func (s responder) internalServerError() {
	s.rw.WriteHeader(http.StatusInternalServerError)
}

func (s responder) created(payload []byte) {
	s.rw.Header().Set("Content-Type", "application/json")
	s.rw.WriteHeader(http.StatusCreated)
	s.rw.Write(payload)
}

func (s responder) badRequest(payload []byte) {
	s.rw.Header().Set("Content-Type", "application/json")
	s.rw.WriteHeader(http.StatusBadRequest)
	s.rw.Write(payload)
}
