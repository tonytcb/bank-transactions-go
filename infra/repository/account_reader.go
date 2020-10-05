package repository

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/tonytcb/bank-transactions-go/domain"
)

// AccountReader exposes account read database operations
type AccountReader struct {
	conn *sql.DB
}

// NewAccountReader build a new AccountReader struct with its dependencies
func NewAccountReader(conn *sql.DB) *AccountReader {
	return &AccountReader{conn: conn}
}

// FindOneByID finds and return one account based in the informed ID
func (a AccountReader) FindOneByID(id *domain.ID) (*domain.Account, error) {
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
