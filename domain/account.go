package domain

import "github.com/pkg/errors"

// AccountRepository represents
type AccountRepository interface {
	Store(*Account) (*ID, error)
}

type Account struct {
	id       *ID
	document *Document
}

func NewAccount(documentNumber DocumentNumber) (*Account, error) {
	document, err := NewDocument(documentNumber)
	if err != nil {
		//return nil, NewErrDomain("document", err.Error()).Wrap("error to create document")
		return nil, err
	}

	return &Account{
		id:       NewID(uint64(0)),
		document: document,
	}, nil
}

func (a *Account) Store(repo AccountRepository) (*Account, error) {
	id, err := repo.Store(a)
	if err != nil {
		return nil, errors.Wrap(err, "error to store account into database")
	}

	return &Account{
		id:       id,
		document: a.document,
	}, nil
}

func (a *Account) Document() *Document {
	return a.document
}

func (a *Account) ID() *ID {
	return a.id
}
