package repository

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/tonytcb/bank-transactions-go/domain"
)

// Transaction exposes transaction database operations
type Transaction struct {
	conn *sql.DB
}

// NewTransaction build a new Transaction struct with its dependencies
func NewTransaction(conn *sql.DB) *Transaction {
	return &Transaction{conn: conn}
}

// Store stores a transaction in the storage
func (t Transaction) Store(transaction *domain.Transaction) (*domain.ID, error) {
	var query = `
		INSERT INTO transactions (account_id, operation_id, amount)
		VALUES (?, ?, ?)
	`

	stmt, err := t.conn.Prepare(query)
	if err != nil {
		return nil, errors.Wrap(err, "prepare statement error")
	}

	result, err := stmt.Exec(transaction.Account().ID().Value(), transaction.Operation().ID().Value(), transaction.Amount())
	if err != nil {
		if v, ok := err.(*mysql.MySQLError); ok {
			return nil, translateMySQLErrors(v)
		}

		return nil, errors.Wrap(err, "unknown database error")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, errors.Wrap(err, "error to read the last inserted id")
	}

	return domain.NewID(uint64(id)), nil
}
