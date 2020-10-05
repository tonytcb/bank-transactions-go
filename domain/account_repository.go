package domain

// AccountRepositoryWriter represents the behaviour of the Account Repository to write operation
type AccountRepositoryWriter interface {
	Store(*Account) (*ID, error)
}

// AccountRepositoryReader represents the behaviour of the Account Repository to read operation
type AccountRepositoryReader interface {
	FindOneByID(*ID) (*Account, error)
}

// AccountRepositoryMock is a fake representation of an AccountRepositoryWriter, useful to create unit tests
type AccountRepositoryMock struct {
	id      *ID
	account *Account
	err     error
}

// NewAccountRepositoryMock builds a new AccountRepositoryMock struct with its mock results
func NewAccountRepositoryMock(id *ID, acc *Account, err error) *AccountRepositoryMock {
	return &AccountRepositoryMock{id: id, account: acc, err: err}
}

// Store stores an account
func (a AccountRepositoryMock) Store(_ *Account) (*ID, error) {
	if a.err != nil {
		return nil, a.err
	}

	return a.id, nil
}

// FindOneByID finds an account by its id
func (a AccountRepositoryMock) FindOneByID(_ *ID) (*Account, error) {
	if a.err != nil {
		return nil, a.err
	}

	return a.account, nil
}
