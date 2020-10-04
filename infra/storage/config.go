package storage

// Config contains all data to create a database connection
type Config struct {
	port     string
	host     string
	password string
	database string
	user     string
}

// NewConfig builds a Config struct
func NewConfig(port, host, password, database, user string) Config {
	return Config{port: port, host: host, password: password, database: database, user: user}
}
