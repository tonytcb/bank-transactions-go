package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// AccountCreator defines the behaviour about how to create an account
type AccountCreator interface {
	Create(string) (int, error)
}

// CreateAccount contains the dependencies to create an account
type CreateAccount struct {
	logger         *log.Logger
	accountCreator AccountCreator
}

// NewCreateAccount creates a new CreateAccount struct with its dependencies
func NewCreateAccount(logger *log.Logger, accountCreator AccountCreator) *CreateAccount {
	return &CreateAccount{logger: logger, accountCreator: accountCreator}
}

// Handler exposes the feature as an http handler
func (h CreateAccount) Handler(rw http.ResponseWriter, req *http.Request) {
	responder := newResponder(rw)

	payload, err := ioutil.ReadAll(req.Body)
	if err != nil {
		h.logger.Println("read payload error:", err)
		responder.internalServerError()
		return
	}

	request := &createAccountPayloadRequest{}
	if err := json.Unmarshal(payload, request); err != nil {
		h.logger.Println("unmarshal payload error:", err)
		responder.internalServerError()
		return
	}
	defer req.Body.Close()

	if errs := request.validate(); errs != nil {
		errResponse := newErrorResponse(errs)

		h.logger.Println("payload doesn't match with the specifications:", errs)
		responder.badRequest(errResponse.Encode())
		return
	}

	id, err := h.accountCreator.Create(request.Document.Number)
	if err != nil {
		errResponse := newErrorResponse(map[string]string{"field-name": err.Error()})

		h.logger.Println("unable to create account:", err)
		responder.badRequest(errResponse.Encode())
		return
	}

	responder.created(createAccountResponse{ID: id}.Encode())
}

type createAccountPayloadRequest struct {
	Document struct {
		Number string `json:"number" validate:"required,number,len=11"`
	}
}

func (c createAccountPayloadRequest) validate() map[string]string {
	if err := validate.Struct(c); err != nil {
		return translateValidations(err.(validator.ValidationErrors))
	}

	return nil
}

type createAccountResponse struct {
	ID int `json:"id"`
}

func (c createAccountResponse) Encode() []byte {
	res, _ := json.Marshal(c)

	return res
}
