package middleware

import (
	"context"
	"net/http"
	"strconv"
	"time"
)

// Logger stores a request identifier in the request context if it was received, otherwise, generate a new one
type RequestID struct {
}

// NewRequestID builds a new Logger struct
func NewRequestID() *RequestID {
	return &RequestID{}
}

// Handler exports RequestID as an http middleware
func (l RequestID) Handler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	requestID := strconv.FormatInt(time.Now().Unix(), 10)

	// if was received the header x-request-id, its value is set in the request context
	if v := r.Header.Values("x-request-id"); len(v) > 0 {
		requestID = v[0]
	}

	ctx := context.WithValue(r.Context(), "request-id", requestID)

	next(w, r.WithContext(ctx))
}
