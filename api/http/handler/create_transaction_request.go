package handler

import "github.com/go-playground/validator/v10"

type createTransactionPayloadRequest struct {
	AccountID   uint64  `json:"account_id" validate:"required,gt=0"`
	OperationID uint64  `json:"operation_id" validate:"required,gt=0"`
	Amount      float64 `json:"amount" validate:"required,gt=0"`
}

// validate returns a map where the key is the field and the value the error description
func (c *createTransactionPayloadRequest) validate() map[string]string {
	if err := validate.Struct(c); err != nil {
		return translateValidations(err.(validator.ValidationErrors))
	}

	return nil
}
