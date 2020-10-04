package domain

// ID represents an valid ID
type ID struct {
	value uint64
}

// NewID build a new ID struct
func NewID(value uint64) *ID {
	return &ID{value: value}
}

// Value returns the valid ID number
func (I ID) Value() uint64 {
	return I.value
}
