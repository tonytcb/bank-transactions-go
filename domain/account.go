package domain

import (
	"time"
)

// Account contains all account's data
type Account struct {
	id        *ID
	document  *Document
	createdAt time.Time
}

// NewAccount creates a new Account struct
func NewAccount(documentNumber DocumentNumber) (*Account, error) {
	document, err := NewDocument(documentNumber)
	if err != nil {
		return nil, err
	}

	return &Account{
		id:       NewID(uint64(0)),
		document: document,
	}, nil
}

// Store stores an account given a Repository
func (a *Account) Store(repo AccountRepository) (*Account, error) {
	id, err := repo.Store(a)
	if err != nil {
		// todo add context to the error returned
		return nil, err
	}

	return &Account{
		id:        id,
		document:  a.document,
		createdAt: time.Now(),
	}, nil
}

// Document returns the document value
func (a *Account) Document() *Document {
	return a.document
}

// ID returns the id value
func (a *Account) ID() *ID {
	return a.id
}

// CreatedAt returns the createdAt value
func (a *Account) CreatedAt() time.Time {
	return a.createdAt
}
