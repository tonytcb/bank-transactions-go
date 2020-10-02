package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/tonytcb/bank-transactions-go/domain"
	"github.com/tonytcb/bank-transactions-go/infra/repository"
)

// AccountCreator defines the behaviour about how to create an account
type AccountCreator interface {
	Create(string) (*domain.Account, error)
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

// Handler exposes the http handler
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
		h.logger.Println("invalid payload:", err)

		errResponse := newErrorResponse(map[string]string{"root": "payload must be a valid JSON"})
		responder.badRequest(errResponse.Encode())

		return
	}
	defer req.Body.Close()

	request.sanitize()

	if errs := request.validate(); errs != nil {
		errResponse := newErrorResponse(errs)

		h.logger.Println("create account payload doesn't match with the specifications:", errs)
		responder.badRequest(errResponse.Encode())
		return
	}

	account, err := h.accountCreator.Create(request.Document.Number)
	if err != nil {
		h.logger.Println("unable to create account:", err)

		if v, ok := err.(*repository.ErrDuplicateEntry); ok {
			h.translateDuplicateError(responder, v)
			return
		}

		if v, ok := err.(*domain.ErrDomain); ok {
			h.translateDomainError(responder, v)
			return
		}

		// unknown error
		responder.internalServerError()
		return
	}

	response := newAccountResponse(
		account.ID().Value(),
		account.Document().Number().String(),
		account.CreatedAt(),
	)

	responder.created(response.Encode())
}

func (h CreateAccount) translateDuplicateError(r *responder, err *repository.ErrDuplicateEntry) {
	errResponse := newErrorResponse(map[string]string{err.Field(): err.Error()})
	r.conflict(errResponse.Encode())
}

func (h CreateAccount) translateDomainError(r *responder, err *domain.ErrDomain) {
	errResponse := newErrorResponse(map[string]string{err.Field(): err.Error()})
	r.unprocessableEntity(errResponse.Encode())
}
