package domain

// TransactionRepositoryWriter represents the behaviour of the Transaction Repository
type TransactionRepositoryWriter interface {
	Store(*Transaction) (*ID, error)
}

// TransactionRepositoryWriterMock is a fake representation of a TransactionRepositoryWriter, useful to create unit tests
type TransactionRepositoryWriterMock struct {
	id  *ID
	err error
}

// NewTransactionRepositoryMock builds a new TransactionRepositoryWriterMock struct with its mock results
func NewTransactionRepositoryMock(id *ID, err error) *TransactionRepositoryWriterMock {
	return &TransactionRepositoryWriterMock{id: id, err: err}
}

// Store stores a transaction
func (t TransactionRepositoryWriterMock) Store(_ *Transaction) (*ID, error) {
	if t.err != nil {
		return nil, t.err
	}

	return t.id, nil
}
