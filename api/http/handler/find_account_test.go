package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/tonytcb/bank-transactions-go/domain"
	"github.com/tonytcb/bank-transactions-go/infra/repository"
)

func TestFindAccount_Handler(t *testing.T) {
	var logger = log.New(fakeWriter{}, "", log.LstdFlags)

	accountOK, _ := domain.NewAccount("00000000191")
	accountOK = accountOK.WithID(domain.NewID(uint64(100))).WithCreateAt(time.Now())

	type fields struct {
		accountFinder AccountFinder
	}

	type args struct {
		id string
	}

	datetimeRegex := `[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}Z`

	tests := []struct {
		name                string
		fields              fields
		args                args
		wantPayloadResponse string
		wantHTTPStatusCode  int
	}{
		// fails
		{
			name: "bad request when the parameter id was not sent",
			fields: fields{
				accountFinder: newFakeAccountFinder(nil, nil),
			},
			args: args{
				id: "x",
			},
			wantPayloadResponse: `{"errors":\[{"field":"id","description":"id must be a valid number"}\]}`,
			wantHTTPStatusCode:  http.StatusBadRequest,
		},
		{
			name: "bad request when the id is zero",
			fields: fields{
				accountFinder: newFakeAccountFinder(nil, nil),
			},
			args: args{
				id: "0",
			},
			wantPayloadResponse: `{"errors":\[{"field":"id","description":"id must be greater than zero"}\]}`,
			wantHTTPStatusCode:  http.StatusBadRequest,
		},
		{
			name: "account not found when the id is not in the storage",
			fields: fields{
				accountFinder: newFakeAccountFinder(nil, repository.NewErrRegisterNotFound("account", "1")),
			},
			args: args{
				id: "1",
			},
			wantPayloadResponse: `{"errors":\[{"field":"id","description":"1 not found"}\]}`,
			wantHTTPStatusCode:  http.StatusNotFound,
		},
		{
			name: "unknown error from account finder",
			fields: fields{
				accountFinder: newFakeAccountFinder(nil, errors.New("some error")),
			},
			args: args{
				id: "100",
			},
			wantPayloadResponse: ``,
			wantHTTPStatusCode:  http.StatusInternalServerError,
		},
		// success
		{
			name: "account found successfully",
			fields: fields{
				accountFinder: newFakeAccountFinder(accountOK, nil),
			},
			args: args{
				id: "100",
			},
			wantPayloadResponse: fmt.Sprintf(`{"id":100,"document":{"number":"00000000191"},"created_at":"%s"}`, datetimeRegex),
			wantHTTPStatusCode:  http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			httpHandler := http.HandlerFunc(NewFindAccount(logger, tt.fields.accountFinder).Handler)
			req, err := http.NewRequest("GET", fmt.Sprintf("/accounts/%s", tt.args.id), nil)
			if err != nil {
				t.Errorf("error to perform GET /accounts/%s request", tt.args.id)
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

type fakeAccountFinder struct {
	account *domain.Account
	err     error
}

func newFakeAccountFinder(account *domain.Account, err error) *fakeAccountFinder {
	return &fakeAccountFinder{account: account, err: err}
}

func (f fakeAccountFinder) Find(*domain.ID) (*domain.Account, error) {
	if f.err != nil {
		return nil, f.err
	}

	return f.account, nil
}
