package usecase

import (
	"github.com/tonytcb/bank-transactions-go/domain"
)

// CreateTransaction contains all the dependencies to create a transaction
type CreateTransaction struct {
	repo domain.TransactionRepository
}

// NewCreateTransaction creates a new CreateTransaction with its dependencies
func NewCreateTransaction(repo domain.TransactionRepository) *CreateTransaction {
	return &CreateTransaction{repo: repo}
}

// Create creates a transaction
func (c CreateTransaction) Create(accountID, operationID *domain.ID, amount float64) (*domain.Transaction, error) {
	transaction, err := domain.NewTransaction(accountID, operationID, amount)
	if err != nil {
		// todo add context to the error
		return nil, err
	}

	t, err := transaction.Store(c.repo)
	if err != nil {
		return nil, err
	}

	return t, nil
}
