package domain

import (
	"strconv"
	"strings"
)

var (
	OperationCompraAVista    = newOperation(uint64(1), "compra a vista")
	OperationCompraParcelada = newOperation(uint64(2), "compra parcelada")
	OperationSaque           = newOperation(uint64(3), "saque")
	OperationPagamento       = newOperation(uint64(4), "Pagamento")

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

func newOperation(id uint64, description string) *Operation {
	return &Operation{id: NewID(id), description: strings.ToUpper(description)}
}

func NewOperation(id *ID) (*Operation, error) {
	if v, ok := operations[id.Value()]; ok {
		return v, nil
	}

	return nil, NewErrDomain("operation", strconv.FormatUint(id.Value(), 10))
}
