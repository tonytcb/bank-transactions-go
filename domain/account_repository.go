package domain

// AccountRepository represents the behaviour of the Account Repository
type AccountRepository interface {
	Store(*Account) (*ID, error)
	FindOneByID(*ID) (*Account, error)
}

// AccountMock is a fake representation of an AccountRepository, useful to create unit tests
type AccountMock struct {
	id      *ID
	account *Account
	err     error
}

// NewAccount builds a new AccountMock struct with its mock results
func NewAccountMock(id *ID, acc *Account, err error) *AccountMock {
	return &AccountMock{id: id, account: acc, err: err}
}

// Store stores an account
func (a AccountMock) Store(_ *Account) (*ID, error) {
	if a.err != nil {
		return nil, a.err
	}

	return a.id, nil
}

// FindOneByID finds an account by its id
func (a AccountMock) FindOneByID(_ *ID) (*Account, error) {
	if a.err != nil {
		return nil, a.err
	}

	return a.account, nil
}
