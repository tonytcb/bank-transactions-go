package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/tonytcb/bank-transactions-go/infra/repository"

	"github.com/tonytcb/bank-transactions-go/domain"
)

// TransactionCreator defines the behaviour about how to create a transaction
type TransactionCreator interface {
	Create(*domain.ID, *domain.ID, float64) (*domain.Transaction, error)
}

// CreateTransaction contains the dependencies to create an transaction
type CreateTransaction struct {
	logger             *log.Logger
	transactionCreator TransactionCreator
}

// NewCreateTransaction creates a new CreateTransaction struct with its dependencies
func NewCreateTransaction(logger *log.Logger, transactionCreator TransactionCreator) *CreateTransaction {
	return &CreateTransaction{logger: logger, transactionCreator: transactionCreator}
}

// Handler exposes the http handler
func (h CreateTransaction) Handler(rw http.ResponseWriter, req *http.Request) {
	responder := newResponder(rw)

	payload, err := ioutil.ReadAll(req.Body)
	if err != nil {
		h.logger.Println("read payload error:", err)
		responder.internalServerError()
		return
	}

	request := &createTransactionPayloadRequest{}
	if err := json.Unmarshal(payload, request); err != nil {
		h.logger.Println("invalid payload:", err)

		errResponse := newErrorResponse(map[string]string{"root": "payload must be a valid JSON"})
		responder.badRequest(errResponse.Encode())

		return
	}
	defer req.Body.Close()

	if errs := request.validate(); errs != nil {
		errResponse := newErrorResponse(errs)

		h.logger.Println("create transaction payload doesn't match with the specifications:", errs)
		responder.badRequest(errResponse.Encode())
		return
	}

	transaction, err := h.transactionCreator.Create(
		domain.NewID(request.AccountID),
		domain.NewID(request.OperationID),
		request.Amount,
	)
	if err != nil {
		h.logger.Println("unable to create transaction:", err)

		if v, ok := err.(*repository.ErrForeignKeyConstraint); ok {
			translateForeignKeyError(responder, v)
			return
		}

		if v, ok := err.(*domain.ErrDomain); ok {
			translateDomainError(responder, v)
			return
		}

		// unknown error
		responder.internalServerError()
		return
	}

	var (
		account   = transaction.Account()
		operation = transaction.Operation()
	)

	response := newTransactionResponse(
		transaction.ID().Value(),
		newAccountResponse(account.ID().Value(), "", account.CreatedAt()),
		newOperationResponse(operation.ID().Value(), operation.Description()),
		transaction.Amount(),
		transaction.CreatedAt(),
	)

	responder.created(response.Encode())
}
