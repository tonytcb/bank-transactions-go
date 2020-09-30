package domain

import "fmt"

type ErrDomain struct {
	field       string
	description string
	rootReason  string
}

func NewErrDomain(field string, description string) *ErrDomain {
	return &ErrDomain{field: field, description: description}
}

func (e ErrDomain) Error() string {
	if e.rootReason != "" {
		return fmt.Sprintf("%s: %s %s", e.rootReason, e.field, e.description)
	}

	return fmt.Sprintf("%s %s", e.field, e.description)
}

func (e ErrDomain) Wrap(rootReason string) *ErrDomain {
	return &ErrDomain{field: e.field, description: e.description, rootReason: rootReason}
}
