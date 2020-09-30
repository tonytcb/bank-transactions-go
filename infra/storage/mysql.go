package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

// NewMySQL creates a new mysql connection
func NewMySQLConnection(c Config) (*sql.DB, error) {
	db, err := sql.Open(c.host, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.user, c.password, c.host, c.port, c.database))
	if err != nil {
		return nil, errors.Wrap(err, "unable to connect to mysql database")
	}

	if err = db.Ping(); err != nil {
		return nil, errors.Wrap(err, "database unavailable")
	}

	return db, nil
}
