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

	"github.com/tonytcb/bank-transactions-go/domain"
	"github.com/tonytcb/bank-transactions-go/infra/repository"
)

func TestCreateAccount_Handler(t *testing.T) {
	var logger = log.New(fakeWriter{}, "", log.LstdFlags)

	accountOK, _ := domain.NewAccount("00000000191")

	datetimeRegex := `[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}Z`

	type fields struct {
		accountCreator AccountCreator
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
			name: "bad request when the payload is empty",
			fields: fields{
				accountCreator: newFakeAccountCreator(nil, nil),
			},
			args: args{
				payload: bytes.NewReader([]byte("")),
			},
			wantPayloadResponse: `{"errors":\[{"field":"root","description":"payload must be a valid JSON"}\]}`,
			wantHTTPStatusCode:  http.StatusBadRequest,
		},
		{
			name: "bad request when the payload has an invalid document number",
			fields: fields{
				accountCreator: newFakeAccountCreator(nil, nil),
			},
			args: args{
				payload: bytes.NewReader([]byte(`{"document": {"number": "000"} }`)),
			},
			wantPayloadResponse: `{"errors":\[{"field":"document.number","description":"number must be 11 characters in length"}\]}`,
			wantHTTPStatusCode:  http.StatusBadRequest,
		},
		{
			name: "internal server error when the payload is corrupted",
			fields: fields{
				accountCreator: newFakeAccountCreator(nil, nil),
			},
			args: args{
				payload: &errReader{},
			},
			wantPayloadResponse: ``,
			wantHTTPStatusCode:  http.StatusInternalServerError,
		},
		{
			name: "internal server error the usecase returns an unknown error",
			fields: fields{
				accountCreator: newFakeAccountCreator(nil, errors.New("unknown error")),
			},
			args: args{
				payload: bytes.NewReader([]byte(`{"document": {"number": "00000000191"} }`)),
			},
			wantPayloadResponse: ``,
			wantHTTPStatusCode:  http.StatusInternalServerError,
		},
		{
			name: "unprocessable entity when the document number is invalid",
			fields: fields{
				accountCreator: newFakeAccountCreator(nil, domain.NewErrDomain("document.number", "'00000000199' is not a valid CPF")),
			},
			args: args{
				payload: bytes.NewReader([]byte(`{"document": {"number": "00000000199"} }`)),
			},
			wantPayloadResponse: `{"errors":\[{"field":"document.number","description":"document.number '00000000199' is not a valid CPF"}\]}`,
			wantHTTPStatusCode:  http.StatusUnprocessableEntity,
		},
		{
			name: "unprocessable entity when the document number is already in storage",
			fields: fields{
				accountCreator: newFakeAccountCreator(nil, repository.NewErrDuplicatedEntry("document_number", "00000000191")),
			},
			args: args{
				payload: bytes.NewReader([]byte(`{"document": {"number": "00000000191"} }`)),
			},
			wantPayloadResponse: `{"errors":\[{"field":"document_number","description":"duplicate entry '00000000191' for field 'document_number'"}\]}`,
			wantHTTPStatusCode:  http.StatusConflict,
		},

		// successes
		{
			name: "account created successfully",
			fields: fields{
				accountCreator: newFakeAccountCreator(accountOK, nil),
			},
			args: args{
				payload: bytes.NewReader([]byte(`{"document": {"number": "00000000191"} }`)),
			},
			wantPayloadResponse: fmt.Sprintf(`{"id":[0-9]+,"document":{"number":"00000000191"},"created_at":"%s"}`, datetimeRegex),
			wantHTTPStatusCode:  http.StatusCreated,
		},
		{
			name: "account created successfully with a formatted document number",
			fields: fields{
				accountCreator: newFakeAccountCreator(accountOK, nil),
			},
			args: args{
				payload: bytes.NewReader([]byte(`{"document": {"number": "000.000.001-91"} }`)),
			},
			wantPayloadResponse: fmt.Sprintf(`{"id":[0-9]+,"document":{"number":"00000000191"},"created_at":"%s"}`, datetimeRegex),
			wantHTTPStatusCode:  http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			httpHandler := http.HandlerFunc(NewCreateAccount(logger, tt.fields.accountCreator).Handler)
			req, err := http.NewRequest("POST", "/accounts", tt.args.payload)
			if err != nil {
				t.Error("error to perform POST /accounts request")
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

			match, err := regexp.MatchString(tt.wantPayloadResponse, gotPayload)
			if err != nil {
				t.Error("Error to validate payload using regex")
			}

			if !match {
				t.Errorf("Payload Response is different from expected, got = %v, want %v", gotPayload, tt.wantPayloadResponse)
				return
			}
		})
	}
}

type fakeWriter struct {
}

func (f fakeWriter) Write(_ []byte) (n int, err error) {
	return 0, nil
}

type fakeAccountCreator struct {
	account *domain.Account
	err     error
}

func newFakeAccountCreator(account *domain.Account, err error) *fakeAccountCreator {
	return &fakeAccountCreator{account: account, err: err}
}

func (f fakeAccountCreator) Create(_ string) (*domain.Account, error) {
	if f.err != nil {
		return nil, f.err
	}

	return f.account, nil
}

type errReader struct {
	io.Reader
}

func (r *errReader) Read(_ []byte) (int, error) {
	return 0, errors.New("test error")
}
