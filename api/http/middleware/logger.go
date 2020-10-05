package middleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// Logger logs the request and response data of the HTTP API
type Logger struct {
	log *log.Logger
}

// NewLogger builds a new Logger struct
func NewLogger(log *log.Logger) *Logger {
	return &Logger{log: log}
}

// Handler exports Logger as an http middleware
func (l Logger) Handler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var (
		start = time.Now()
		rec   = newStatusRecorder(w)
	)

	l.log.Println(l.requestData(r))

	next.ServeHTTP(&rec, r)

	l.log.Println(rec.responseData(start))
}

func (l Logger) requestData(r *http.Request) string {
	payload, _ := getPayload(r)

	v, err := json.Marshal(map[string]interface{}{
		"request": map[string]interface{}{
			"http_method": r.Method,
			"path":        r.URL.Path,
			"headers":     r.Header,
			"payload":     compactPayload(payload),
		},
	})

	if err != nil {
		return ""
	}

	return string(v)
}

func getPayload(r *http.Request) (string, error) {
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	// fill back the body buffer
	r.Body = ioutil.NopCloser(bytes.NewBuffer(payload))

	return string(payload), nil
}

func compactPayload(payload string) string {
	var (
		buffer    = new(bytes.Buffer)
		jsonBytes = []byte(payload)
	)

	if err := json.Compact(buffer, jsonBytes); err != nil {
		return string(jsonBytes)
	}

	// remove line breaks
	return strings.Replace(payload, "\n", "", -1)
}

// statusRecorder armazena informações da resposta da API HTTP, sobrescrevendo o http.responseWriter
type statusRecorder struct {
	http.ResponseWriter
	status int
	body   []byte
}

// default response
func newStatusRecorder(w http.ResponseWriter) statusRecorder {
	return statusRecorder{ResponseWriter: w, status: http.StatusOK, body: []byte("")}
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code

	rec.ResponseWriter.WriteHeader(code)
}

func (rec *statusRecorder) Write(body []byte) (int, error) {
	rec.body = body

	_, err := rec.ResponseWriter.Write(body)

	return rec.status, err
}

func (rec *statusRecorder) responseData(start time.Time) string {
	var payload = string(rec.body)
	if payload == "" {
		payload = `{}`
	}

	v, err := json.Marshal(map[string]interface{}{
		"response": map[string]interface{}{
			"http_status": rec.status,
			"payload":     payload,
			"time":        time.Since(start).Seconds(),
		},
	})

	if err != nil {
		return ""
	}

	return string(v)
}
