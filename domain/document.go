package domain

import (
	"fmt"

	"github.com/Nhanderu/brdoc"
)

// DocumentNumber represents the document's number
type DocumentNumber string

// String cast the document number value to string
func (d DocumentNumber) String() string {
	return string(d)
}

// Document represents a customer's document
type Document struct {
	number DocumentNumber
}

// NewDocument build a new Documents with its dependencies
func NewDocument(number DocumentNumber) (*Document, error) {
	if brdoc.IsCPF(string(number)) {
		return &Document{number: number}, nil
	}

	return nil, NewErrDomain("document.number", fmt.Sprintf("'%s' is not a valid document number", number))
}

// Number returns the value of document number
func (d Document) Number() DocumentNumber {
	return d.number
}
