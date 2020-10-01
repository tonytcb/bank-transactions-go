package repository

import (
	"database/sql"
	"regexp"

	"github.com/go-sql-driver/mysql"
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

// Store stores an account in the storage
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

		if v, ok := err.(*mysql.MySQLError); ok {
			return nil, a.translateMySqlErrors(v)
		}

		return nil, errors.Wrap(err, "database error")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, errors.Wrap(err, "error to read the last inserted id")
	}

	return domain.NewID(uint64(id)), nil
}

func (a Account) translateMySqlErrors(err *mysql.MySQLError) error {
	const (
		duplicateEntryCode = 1062
	)

	if err.Number == duplicateEntryCode {
		duplicatedErrorRegex := regexp.MustCompile(`Duplicate entry '(\w+)' for key '(\w+)'`)

		if match := duplicatedErrorRegex.FindAllStringSubmatch(err.Message, -1); len(match) > 0 {
			return NewErrDuplicatedEntry(match[0][2], match[0][1])
		}
	}

	return err
}
