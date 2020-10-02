package handler

import "net/http"

type responder struct {
	rw http.ResponseWriter
}

func newResponder(rw http.ResponseWriter) *responder {
	return &responder{rw: rw}
}

func (s responder) created(payload []byte) {
	s.rw.Header().Set("Content-Type", "application/json")
	s.rw.WriteHeader(http.StatusCreated)
	s.rw.Write(payload)
}

func (s responder) ok(payload []byte) {
	s.rw.Header().Set("Content-Type", "application/json")
	s.rw.WriteHeader(http.StatusOK)
	s.rw.Write(payload)
}

func (s responder) internalServerError() {
	s.rw.WriteHeader(http.StatusInternalServerError)
}

func (s responder) badRequest(payload []byte) {
	s.rw.Header().Set("Content-Type", "application/json")
	s.rw.WriteHeader(http.StatusBadRequest)
	s.rw.Write(payload)
}

func (s responder) notFound(payload []byte) {
	s.rw.Header().Set("Content-Type", "application/json")
	s.rw.WriteHeader(http.StatusNotFound)
	s.rw.Write(payload)
}

func (s responder) conflict(payload []byte) {
	s.rw.Header().Set("Content-Type", "application/json")
	s.rw.WriteHeader(http.StatusConflict)
	s.rw.Write(payload)
}

func (s responder) unprocessableEntity(payload []byte) {
	s.rw.Header().Set("Content-Type", "application/json")
	s.rw.WriteHeader(http.StatusUnprocessableEntity)
	s.rw.Write(payload)
}
