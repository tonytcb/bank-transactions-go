package handler

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

type createAccountPayloadRequest struct {
	Document struct {
		Number string `json:"number" validate:"required,number,len=11"`
	}
}

func (c *createAccountPayloadRequest) sanitize() {
	numberRegex := regexp.MustCompile(`[0-9]+`)

	if parts := numberRegex.FindAllString(c.Document.Number, -1); len(parts) > 0 {
		c.Document.Number = strings.Join(parts, "")
	}
}

func (c *createAccountPayloadRequest) validate() map[string]string {
	if err := validate.Struct(c); err != nil {
		return translateValidations(err.(validator.ValidationErrors))
	}

	return nil
}
