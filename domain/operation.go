package domain

import (
	"fmt"
	"strings"
)

var (
	// OperationCompraAVista representa uma operação de compra a vista
	OperationCompraAVista = newOperation(uint64(1), "compra a vista")

	// OperationCompraParcelada representa uma operação de compra parcelada
	OperationCompraParcelada = newOperation(uint64(2), "compra parcelada")

	// OperationSaque representa uma operação de saque
	OperationSaque = newOperation(uint64(3), "saque")

	// OperationPagamento representa uma operação de pagamento
	OperationPagamento = newOperation(uint64(4), "Pagamento")

	operations = map[uint64]*Operation{
		OperationCompraAVista.id.Value():    OperationCompraAVista,
		OperationCompraParcelada.id.Value(): OperationCompraParcelada,
		OperationSaque.id.Value():           OperationSaque,
		OperationPagamento.id.Value():       OperationPagamento,
	}
)

// Operation contains all data to recognize the type of a transaction
type Operation struct {
	id          *ID
	description string
}

// IsIncoming checks if the operation is an incoming operation
func (o Operation) IsIncoming() bool {
	ids := []uint64{uint64(4)}

	for _, v := range ids {
		if v == o.id.Value() {
			return true
		}
	}

	return false
}

// ID returns the id value
func (o Operation) ID() *ID {
	return o.id
}

// Description returns the description value
func (o Operation) Description() string {
	return o.description
}

func newOperation(id uint64, description string) *Operation {
	return &Operation{id: NewID(id), description: strings.ToUpper(description)}
}

// NewOperation creates a valid Operation struct
func NewOperation(id *ID) (*Operation, error) {
	if v, ok := operations[id.Value()]; ok {
		return v, nil
	}

	description := fmt.Sprintf("'%d' is not a valid operation id", id.Value())

	return nil, NewErrDomain("operation", description)
}
