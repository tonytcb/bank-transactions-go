package domain

// TransactionRepository represents the behaviour of the Transaction Repository
type TransactionRepository interface {
	Store(*Transaction) (*ID, error)
}

// TransactionRepositoryMock is a fake representation of a TransactionRepository, useful to create unit tests
type TransactionRepositoryMock struct {
	id  *ID
	err error
}

// NewTransactionRepositoryMock builds a new TransactionRepositoryMock struct with its mock results
func NewTransactionRepositoryMock(id *ID, err error) *TransactionRepositoryMock {
	return &TransactionRepositoryMock{id: id, err: err}
}

// Store stores a transaction
func (t TransactionRepositoryMock) Store(_ *Transaction) (*ID, error) {
	if t.err != nil {
		return nil, t.err
	}

	return t.id, nil
}
