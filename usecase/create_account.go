package usecase

import (
	"github.com/tonytcb/bank-transactions-go/domain"
)

// CreateAccount contains all the dependencies to create an account
type CreateAccount struct {
	repo domain.AccountRepository
}

// NewCreateAccount creates a new CreateAccount with its dependencies
func NewCreateAccount(repo domain.AccountRepository) *CreateAccount {
	return &CreateAccount{repo: repo}
}

// Create creates a account
func (c CreateAccount) Create(documentNumber string) (*domain.Account, error) {
	account, err := domain.NewAccount(domain.DocumentNumber(documentNumber))
	if err != nil {
		// todo add context to the error
		return nil, err
	}

	acc, err := account.Store(c.repo)
	if err != nil {
		return nil, err
	}

	return acc, nil
}
