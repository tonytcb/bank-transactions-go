package handler

import (
	"encoding/json"
	"time"
)

type operationResponse struct {
	ID          uint64 `json:"id"`
	Description string `json:"type"`
}

func newOperationResponse(ID uint64, description string) operationResponse {
	return operationResponse{ID: ID, Description: description}
}

type transactionResponse struct {
	ID        uint64            `json:"id"`
	Account   accountResponse   `json:"account,omitempty"`
	Operation operationResponse `json:"operation"`
	Amount    float64           `json:"amount"`
	CreatedAt string            `json:"created_at"`
}

func newTransactionResponse(id uint64, acc accountResponse, op operationResponse, amount float64, t time.Time) transactionResponse {
	return transactionResponse{
		ID:        id,
		Account:   acc,
		Operation: op,
		Amount:    amount,
		CreatedAt: t.UTC().Format(time.RFC3339),
	}
}

func (c transactionResponse) Encode() []byte {
	res, _ := json.Marshal(c)

	return res
}
