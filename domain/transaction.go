package domain

import (
	"time"
)

// Transaction represents a Transaction in the Domain
type Transaction struct {
	id        *ID
	account   *Account
	operation *Operation
	amount    float64
	createdAt time.Time
}

// NewTransaction builds a new Transaction struct
func NewTransaction(accountID *ID, operationID *ID, amount float64) (*Transaction, error) {
	operation, err := NewOperation(operationID)
	if err != nil {
		return nil, err
	}

	account := &Account{id: accountID}

	if !operation.IsIncoming() {
		amount = -amount
	}

	return &Transaction{
		id:        NewID(0),
		account:   account,
		operation: operation,
		amount:    amount,
	}, nil
}

// Store stores a transaction given a repository
func (t *Transaction) Store(repo TransactionRepository) (*Transaction, error) {
	id, err := repo.Store(t)
	if err != nil {
		return nil, err
	}

	return &Transaction{
		id:        id,
		account:   t.Account(),
		operation: t.Operation(),
		amount:    t.Amount(),
		createdAt: time.Now(),
	}, nil
}

// ID returns the transaction's id
func (t *Transaction) ID() *ID {
	return t.id
}

// Account returns the account used to register the transaction
func (t *Transaction) Account() *Account {
	return t.account
}

// Operation returns the type of then operation of the transaction
func (t *Transaction) Operation() *Operation {
	return t.operation
}

// Amount returns the amount
func (t *Transaction) Amount() float64 {
	return t.amount
}

// CreatedAt returns the createdAt value
func (t *Transaction) CreatedAt() time.Time {
	return t.createdAt
}

// WithAccount returns a new Transaction struct with the informed account
func (t *Transaction) WithAccount(a *Account) *Transaction {
	return &Transaction{
		account:   a,
		id:        t.ID(),
		operation: t.Operation(),
		amount:    t.Amount(),
		createdAt: t.CreatedAt(),
	}
}

// WithID returns a new Transaction struct with the informed account
func (t *Transaction) WithID(id *ID) *Transaction {
	return &Transaction{
		id:        id,
		account:   t.Account(),
		operation: t.Operation(),
		amount:    t.Amount(),
		createdAt: t.CreatedAt(),
	}
}
