package domain

import "fmt"

type ErrDomain struct {
	field       string
	description string
}

func (e ErrDomain) Field() string {
	return e.field
}

func (e ErrDomain) Description() string {
	return e.description
}

func NewErrDomain(field string, description string) *ErrDomain {
	return &ErrDomain{field: field, description: description}
}

func (e ErrDomain) Error() string {
	return fmt.Sprintf("%s %s", e.Field(), e.Description())
}
