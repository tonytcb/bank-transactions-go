package domain

import "fmt"

// ErrDomain represents an well known domain error
type ErrDomain struct {
	field       string
	description string
}

// Field returns the field where occurred the error
func (e ErrDomain) Field() string {
	return e.field
}

// Description returns the description of the error
func (e ErrDomain) Description() string {
	return e.description
}

// NewErrDomain build a new ErrDomain struct
func NewErrDomain(field string, description string) *ErrDomain {
	return &ErrDomain{field: field, description: description}
}

// Error returns a formatted error message
func (e ErrDomain) Error() string {
	return fmt.Sprintf("%s %s", e.Field(), e.Description())
}
