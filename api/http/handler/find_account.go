package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/tonytcb/bank-transactions-go/domain"
	"github.com/tonytcb/bank-transactions-go/infra/repository"
)

// AccountFinder defines the behaviour about how to find an account
type AccountFinder interface {
	Find(*domain.ID) (*domain.Account, error)
}

// CreateAccount contains the dependencies to find an account
type FindAccount struct {
	logger        *log.Logger
	accountFinder AccountFinder
}

// NewFindAccount creates a new FindAccount struct
func NewFindAccount(logger *log.Logger, accountFinder AccountFinder) *FindAccount {
	return &FindAccount{logger: logger, accountFinder: accountFinder}
}

// Handler exposes the http handler
func (f FindAccount) Handler(rw http.ResponseWriter, req *http.Request) {
	responder := newResponder(rw)

	idParam, err := f.extractParamGetID(req)
	if err != nil {
		f.logger.Println("invalid account id:", err)

		errResponse := newErrorResponse(map[string]string{"id": err.Error()})
		responder.badRequest(errResponse.Encode())
		return
	}

	account, err := f.accountFinder.Find(domain.NewID(idParam))
	if err != nil {
		if _, ok := err.(*repository.ErrRegisterNotFound); ok {
			f.logger.Println("account not found:", err)
			errResponse := newErrorResponse(map[string]string{"id": fmt.Sprintf("%d not found", idParam)})
			responder.notFound(errResponse.Encode())
			return
		}

		f.logger.Println("unknown error:", err)
		responder.internalServerError()
		return
	}

	response := newAccountResponse(
		account.ID().Value(),
		account.Document().Number().String(),
		account.CreatedAt(),
	)

	f.logger.Println("account found:", string(response.Encode()))

	responder.ok(response.Encode())
}

func (f FindAccount) extractParamGetID(req *http.Request) (uint64, error) {
	const position = 2

	p := strings.Split(req.URL.Path, "/")

	if len(p) < (position + 1) {
		return 0, errors.New("parameter id not found")
	}

	id, err := strconv.Atoi(p[position])
	if err != nil {
		return 0, errors.New("id must be a valid number")
	}

	if id <= 0 {
		return 0, errors.New("id must be greater than zero")
	}

	return uint64(id), nil
}
