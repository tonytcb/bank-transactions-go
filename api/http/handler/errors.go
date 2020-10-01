package handler

import "encoding/json"

type err struct {
	Field       string `json:"field"`
	Description string `json:"description"`
}

type errorResponse struct {
	Errors []err `json:"errors"`
}

func newErrorResponse(e map[string]string) *errorResponse {
	var errs []err
	for i, v := range e {
		errs = append(errs, err{Field: i, Description: v})
	}
	return &errorResponse{Errors: errs}
}

func (e errorResponse) Encode() []byte {
	res, _ := json.Marshal(e)

	return res
}
