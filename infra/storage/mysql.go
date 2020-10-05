package storage

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" // Register MySQL operations
	"github.com/pkg/errors"
)

// NewMySQLConnection creates a new mysql connection
func NewMySQLConnection(c Config) (*sql.DB, error) {
	toRetry := func() (*sql.DB, error) {
		db, err := sql.Open(c.host, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.user, c.password, c.host, c.port, c.database))
		if err != nil {
			return nil, errors.Wrap(err, "unable to connect to mysql database")
		}

		if err = db.Ping(); err != nil {
			return nil, errors.Wrap(err, "database unavailable")
		}

		return db, nil
	}

	return retry(toRetry, 20)
}

func retry(fn func() (*sql.DB, error), maxRetries int) (*sql.DB, error) {
	i := 0
	for {
		conn, err := fn()
		if err == nil {
			return conn, nil
		}

		if i >= maxRetries {
			return nil, errors.Wrap(err, fmt.Sprintf("error after %d attemps", i))
		}

		i++
		time.Sleep(1 * time.Second)
	}
}
