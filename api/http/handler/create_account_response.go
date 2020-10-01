package handler

import (
	"encoding/json"
	"time"
)

type createAccountResponseDocumentResponse struct {
	Number string `json:"number"`
}

type createAccountResponse struct {
	ID        uint64                                `json:"id"`
	Document  createAccountResponseDocumentResponse `json:"document"`
	CreatedAt string                                `json:"created_at"`
}

func newCreateAccountResponse(ID uint64, documentNumber string, createdAt time.Time) createAccountResponse {
	return createAccountResponse{
		ID:        ID,
		Document:  createAccountResponseDocumentResponse{Number: documentNumber},
		CreatedAt: createdAt.Format(time.RFC3339),
	}
}

func (c createAccountResponse) Encode() []byte {
	res, _ := json.Marshal(c)

	return res
}
