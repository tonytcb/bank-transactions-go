package usecase

import (
	"reflect"
	"testing"

	"github.com/tonytcb/bank-transactions-go/infra/repository"

	"github.com/tonytcb/bank-transactions-go/domain"
)

func TestCreateAccount(t *testing.T) {
	type fields struct {
		repo domain.AccountRepository
	}
	type args struct {
		documentNumber string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "domain error when the document number is invalid",
			fields: fields{
				repo: domain.NewAccountMock(nil, nil, nil),
			},
			args: args{
				documentNumber: "00000000000",
			},
			wantErr: domain.NewErrDomain("document.number", "'00000000000' is not a valid CPF"),
		},
		{
			name: "repository error when the document numbers is duplicate",
			fields: fields{
				repo: domain.NewAccountMock(nil, nil, repository.NewErrDuplicatedEntry("document number", "duplicate entry 00000000191")),
			},
			args: args{
				documentNumber: "00000000191",
			},
			wantErr: repository.NewErrDuplicatedEntry("document number", "duplicate entry 00000000191"),
		},
		{
			name: "account created successfully",
			fields: fields{
				repo: domain.NewAccountMock(domain.NewID(1000), nil, nil),
			},
			args: args{
				documentNumber: "00000000191",
			},
			wantErr: repository.NewErrDuplicatedEntry("document number", "duplicate entry 00000000191"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCreateAccount(tt.fields.repo)

			got, err := c.Create(tt.args.documentNumber)
			if (err != nil) && !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			if got.ID().Value() <= 0 {
				t.Errorf("Invalid Account result: ID it must be greater than zero")
				return
			}
		})
	}
}
