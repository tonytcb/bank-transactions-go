package handler

import (
	"encoding/json"
	"time"
)

type documentResponse struct {
	Number string `json:"number,omitempty"`
}

type accountResponse struct {
	ID        uint64           `json:"id,omitempty"`
	Document  documentResponse `json:"document,omitempty"`
	CreatedAt string           `json:"created_at,omitempty"`
}

func newAccountResponse(ID uint64, documentNumber string, createdAt time.Time) accountResponse {
	t := createdAt.UTC().Format(time.RFC3339)
	if createdAt.Year() == 1 {
		t = ""
	}

	return accountResponse{
		ID:        ID,
		Document:  documentResponse{Number: documentNumber},
		CreatedAt: t,
	}
}

func (c accountResponse) Encode() []byte {
	res, _ := json.Marshal(c)

	return res
}
