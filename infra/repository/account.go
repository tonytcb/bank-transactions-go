package repository

import (
	"database/sql"
	"strconv"
	"time"

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

// FindOneByID finds and return one account based in the informed ID
func (a Account) FindOneByID(id *domain.ID) (*domain.Account, error) {
	var (
		documentNumber     string
		createdAtTimestamp []uint8
		query              = `SELECT document_number, created_at FROM accounts WHERE id = ?`
	)

	if err := a.conn.QueryRow(query, id.Value()).Scan(&documentNumber, &createdAtTimestamp); err != nil {
		if err == sql.ErrNoRows {
			return nil, NewErrRegisterNotFound("id", strconv.FormatUint(id.Value(), 10))
		}

		return nil, err
	}

	createdAt, err := timestampToTime(createdAtTimestamp)
	if err != nil {
		createdAt = time.Time{}
	}

	account, err := domain.NewAccount(domain.DocumentNumber(documentNumber))
	if err != nil {
		return nil, NewErrLoadInvalidData("accounts")
	}

	return account.WithID(id).WithCreateAt(createdAt), nil
}

func timestampToTime(t []uint8) (time.Time, error) {
	parsedTime, err := time.Parse("2006-01-02 15:04:05", string(t))
	if err != nil {
		return time.Time{}, err
	}

	return parsedTime, nil
}
