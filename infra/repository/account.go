package repository

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/tonytcb/bank-transactions-go/domain"
)

// Account exposes account database operations
type Account struct {
	conn *sql.DB
}

// NewAccount build a new Account struct with its dependencies
func NewAccount(conn *sql.DB) *Account {
	return &Account{conn: conn}
}

// Store stores an account into the database
func (a Account) Store(acc *domain.Account) (*domain.ID, error) {
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
		return nil, errors.Wrap(err, "error to create account")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, errors.Wrap(err, "error to read the last inserted id")
	}

	return domain.NewID(uint64(id)), nil
}
