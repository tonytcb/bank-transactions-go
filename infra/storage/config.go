package storage

type Config struct {
	port     string
	host     string
	password string
	database string
	user     string
}

func NewConfig(port, host, password, database, user string) Config {
	return Config{port: port, host: host, password: password, database: database, user: user}
}
