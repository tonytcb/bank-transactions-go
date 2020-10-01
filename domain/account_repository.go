package domain

// AccountRepository represents the behaviour of the Account Repository
type AccountRepository interface {
	Store(*Account) (*ID, error)
}

// AccountMock represents a fake representation of an Account
type AccountMock struct {
	result *ID
	err    error
}

// NewAccount build a new AccountMock struct with its mock results
func NewAccountMock(result *ID, err error) *AccountMock {
	return &AccountMock{result: result, err: err}
}

// Store stores an account
func (a AccountMock) Store(_ *Account) (*ID, error) {
	if a.err != nil {
		return nil, a.err
	}

	return a.result, nil
}
