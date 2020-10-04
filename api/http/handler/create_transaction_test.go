package handler

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/tonytcb/bank-transactions-go/infra/repository"

	"github.com/tonytcb/bank-transactions-go/domain"
)

func TestCreateTransaction_Handler(t *testing.T) {
	var (
		logger                 = log.New(fakeWriter{}, "", log.LstdFlags)
		foreignKeyAccountError = repository.NewErrForeignKeyConstraint("accounts", "accountfk1", "account_id", "id")
		operationError         = domain.NewErrDomain("operation", "'10' is not a valid operation id")
		datetimeRegex          = `[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}Z`
	)

	transactionOK, _ := domain.NewTransaction(domain.NewID(1), domain.NewID(4), 100)

	type fields struct {
		transactionCreator TransactionCreator
	}
	type args struct {
		payload io.Reader
	}
	tests := []struct {
		name                string
		fields              fields
		args                args
		wantPayloadResponse string
		wantHTTPStatusCode  int
	}{
		// fails
		{
			name: "internal server error when the payload is corrupted",
			fields: fields{
				transactionCreator: newFakeTransactionCreator(nil, nil),
			},
			args: args{
				payload: &errReader{},
			},
			wantPayloadResponse: ``,
			wantHTTPStatusCode:  http.StatusInternalServerError,
		},
		{
			name: "bad request when the payload is empty",
			fields: fields{
				transactionCreator: newFakeTransactionCreator(nil, nil),
			},
			args: args{
				payload: bytes.NewReader([]byte("")),
			},
			wantPayloadResponse: `{"errors":\[{"field":"root","description":"payload must be a valid JSON"}\]}`,
			wantHTTPStatusCode:  http.StatusBadRequest,
		},
		{
			name: "bad request when the payload has an invalid amount",
			fields: fields{
				transactionCreator: newFakeTransactionCreator(nil, nil),
			},
			args: args{
				payload: bytes.NewReader([]byte(`{"account_id": 1, "operation_id": 1, "amount": -100.00}`)),
			},
			wantPayloadResponse: `{"errors":\[{"field":"amount","description":"amount must be greater than 0"}\]}`,
			wantHTTPStatusCode:  http.StatusBadRequest,
		},
		{
			name: "bad request when the payload has an invalid amount",
			fields: fields{
				transactionCreator: newFakeTransactionCreator(nil, nil),
			},
			args: args{
				payload: bytes.NewReader([]byte(`{"account_id": 1, "operation_id": 1, "amount": -100.00}`)),
			},
			wantPayloadResponse: `{"errors":\[{"field":"amount","description":"amount must be greater than 0"}\]}`,
			wantHTTPStatusCode:  http.StatusBadRequest,
		},
		{
			name: "bad request when the payload has not an account_id",
			fields: fields{
				transactionCreator: newFakeTransactionCreator(nil, nil),
			},
			args: args{
				payload: bytes.NewReader([]byte(`{"operation_id": 1, "amount": 100.00}`)),
			},
			wantPayloadResponse: `{"errors":\[{"field":"account_id","description":"account_id is a required field"}\]}`,
			wantHTTPStatusCode:  http.StatusBadRequest,
		},
		{
			name: "bad request when the the account was not found",
			fields: fields{
				transactionCreator: newFakeTransactionCreator(nil, foreignKeyAccountError),
			},
			args: args{
				payload: bytes.NewReader([]byte(`{"account_id": 101, "operation_id": 1, "amount": 100.00}`)),
			},
			wantPayloadResponse: `{"errors":\[{"field":"account_id","description":"'account_id' not found"}\]}`,
			wantHTTPStatusCode:  http.StatusUnprocessableEntity,
		},
		{
			name: "bad request when the the operation is invalid",
			fields: fields{
				transactionCreator: newFakeTransactionCreator(nil, operationError),
			},
			args: args{
				payload: bytes.NewReader([]byte(`{"account_id": 1, "operation_id": 10, "amount": 100.00}`)),
			},
			wantPayloadResponse: `{"errors":\[{"field":"operation","description":"operation '10' is not a valid operation id"}\]}`,
			wantHTTPStatusCode:  http.StatusUnprocessableEntity,
		},
		{
			name: "internal server error when returns an unknown error",
			fields: fields{
				transactionCreator: newFakeTransactionCreator(nil, errors.New("unknown error")),
			},
			args: args{
				payload: bytes.NewReader([]byte(`{"account_id": 1, "operation_id": 4, "amount": 100.00}`)),
			},
			wantPayloadResponse: ``,
			wantHTTPStatusCode:  http.StatusInternalServerError,
		},
		{
			name: "transaction created successfully",
			fields: fields{
				transactionCreator: newFakeTransactionCreator(transactionOK.WithID(domain.NewID(50)), nil),
			},
			args: args{
				payload: bytes.NewReader([]byte(`{"account_id": 1, "operation_id": 4, "amount": 100.00}`)),
			},
			wantPayloadResponse: fmt.Sprintf(`{"id":50,"account":{"id":1,"document":{}},"operation":{"id":4,"type":"PAGAMENTO"},"amount":100,"created_at":"%s"}`, datetimeRegex),
			wantHTTPStatusCode:  http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			httpHandler := http.HandlerFunc(NewCreateTransaction(logger, tt.fields.transactionCreator).Handler)
			req, err := http.NewRequest("POST", "/transactions", tt.args.payload)
			if err != nil {
				t.Error("error to perform POST /transactions request")
			}

			httpHandler.ServeHTTP(rr, req)

			var (
				gotHTTPStatusCode = rr.Code
				gotPayload        = rr.Body.String()
			)

			if gotHTTPStatusCode != tt.wantHTTPStatusCode {
				t.Errorf("HTTP Status Code is different from expected, got = %v, want %v", gotHTTPStatusCode, tt.wantHTTPStatusCode)
				return
			}

			match := regexp.MustCompile(tt.wantPayloadResponse).MatchString(gotPayload)
			if !match {
				t.Errorf("Payload Response is different from expected, got = %v, want %v", gotPayload, tt.wantPayloadResponse)
				return
			}
		})
	}
}

type fakeTransactionCreator struct {
	transaction *domain.Transaction
	err         error
}

func newFakeTransactionCreator(transaction *domain.Transaction, err error) *fakeTransactionCreator {
	return &fakeTransactionCreator{transaction: transaction, err: err}
}

func (f fakeTransactionCreator) Create(*domain.ID, *domain.ID, float64) (*domain.Transaction, error) {
	if f.err != nil {
		return nil, f.err
	}

	return f.transaction, nil
}
