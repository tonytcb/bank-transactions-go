package repository

import "fmt"

// ErrDuplicateEntry represents a duplicate entry error
type ErrDuplicateEntry struct {
	field string
	value string
}

// NewErrDuplicatedEntry creates an ErrDuplicateEntry struct
func NewErrDuplicatedEntry(field string, value string) *ErrDuplicateEntry {
	return &ErrDuplicateEntry{field: field, value: value}
}

// Error returns the formatted error message
func (e ErrDuplicateEntry) Error() string {
	return fmt.Sprintf(`duplicate entry '%s' for field '%s'`, e.value, e.field)
}

func (e ErrDuplicateEntry) Value() string {
	return e.value
}

func (e ErrDuplicateEntry) Field() string {
	return e.field
}
