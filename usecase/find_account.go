package usecase

import (
	"github.com/tonytcb/bank-transactions-go/domain"
)

// FindAccount contains all the dependencies to find an account
type FindAccount struct {
	repo domain.AccountRepositoryReader
}

// NewFindAccount creates a new FindAccount with its dependencies
func NewFindAccount(repo domain.AccountRepositoryReader) *FindAccount {
	return &FindAccount{repo: repo}
}

// Find finds an account by its id
func (f FindAccount) Find(id *domain.ID) (*domain.Account, error) {
	account, err := f.repo.FindOneByID(id)
	if err != nil {
		return nil, err
	}

	return account, nil
}
