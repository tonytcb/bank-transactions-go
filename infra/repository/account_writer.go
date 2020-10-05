package repository

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/tonytcb/bank-transactions-go/domain"
)

// AccountWriter exposes account write database operations
type AccountWriter struct {
	conn *sql.DB
}

// NewAccountWriter build a new AccountWriter struct with its dependencies
func NewAccountWriter(conn *sql.DB) *AccountWriter {
	return &AccountWriter{conn: conn}
}

// Store stores an account in the storage
func (a AccountWriter) Store(acc *domain.Account) (*domain.ID, error) {
	var query = `
		INSERT INTO accounts (document_number)
		VALUES (?)
	`

	stmt, err := a.conn.Prepare(query)
	if err != nil {
		return nil, errors.Wrap(err, "prepare statement error")
	}

	result, err := stmt.Exec(acc.Document().Number().String())
	if err != nil {
		if v, ok := err.(*mysql.MySQLError); ok {
			return nil, translateMySQLErrors(v)
		}

		return nil, errors.Wrap(err, "database error")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, errors.Wrap(err, "error to read the last inserted id")
	}

	return domain.NewID(uint64(id)), nil
}
