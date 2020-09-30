package api

// Server defines the behaviour of the server, regardless the protocol (http, amqp, ...) used in the implementation
type Server interface {
	Listen()
}