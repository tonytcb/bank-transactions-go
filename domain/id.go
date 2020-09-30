package domain

type ID struct {
	value uint64
}

func NewID(value uint64) *ID {
	return &ID{value: value}
}

func (I ID) Value() uint64 {
	return I.value
}
