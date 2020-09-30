package domain

import (
	"errors"

	"github.com/Nhanderu/brdoc"
)

// DocumentNumber represents the document's number
type DocumentNumber string

func (d DocumentNumber) String() string {
	return string(d)
}

// Document represents a customer's document
type Document struct {
	number DocumentNumber
}

// NewDocument build a new Documents with its dependencies
func NewDocument(number DocumentNumber) (*Document, error) {
	if !brdoc.IsCPF(string(number)) {
		//return nil, NewErrDomain("number", "must be a valid brazilian document number")
		return nil, errors.New("document number must be a valid CPF")
	}

	return &Document{number: number}, nil
}

// Number returns the value of document number
func (d Document) Number() DocumentNumber {
	return d.number
}
