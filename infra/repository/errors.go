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

// Field returns the duplicate value
func (e ErrDuplicateEntry) Value() string {
	return e.value
}

// Field returns the field value
func (e ErrDuplicateEntry) Field() string {
	return e.field
}

// --

// ErrRegisterNotFound represents an error when the register was not found
type ErrRegisterNotFound struct {
	field string
	value string
}

// NewErrRegisterNotFound creates an ErrDuplicateEntry struct
func NewErrRegisterNotFound(field, value string) *ErrRegisterNotFound {
	return &ErrRegisterNotFound{field: field, value: value}
}

// Field returns the field not found
func (e ErrRegisterNotFound) Field() string {
	return e.field
}

// Field returns the duplicate value
func (e ErrRegisterNotFound) Value() string {
	return e.value
}

// Error returns the formatted error message
func (e ErrRegisterNotFound) Error() string {
	return fmt.Sprintf(`'%s' '%s' not found`, e.field, e.value)
}

// --

// ErrRegisterNotFound represents an error when occurred an error to load data
type ErrLoadInvalidData struct {
	field string
}

// NewErrLoadInvalidData creates an ErrLoadInvalidData struct
func NewErrLoadInvalidData(field string) *ErrLoadInvalidData {
	return &ErrLoadInvalidData{field: field}
}

// Error returns the formatted error message
func (e ErrLoadInvalidData) Error() string {
	return fmt.Sprintf(`error to load values from '%s'`, e.field)
}
