package usecase

import (
	"errors"
	"reflect"
	"testing"

	"github.com/tonytcb/bank-transactions-go/domain"
)

func TestCreateTransaction_Create(t *testing.T) {
	transaction := &domain.Transaction{}

	type fields struct {
		repo domain.TransactionRepositoryWriter
	}
	type args struct {
		accountID   *domain.ID
		operationID *domain.ID
		amount      float64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Transaction
		wantErr error
	}{
		{
			name: "domain error when the operation is not valid",
			fields: fields{
				repo: domain.NewTransactionRepositoryMock(nil, nil),
			},
			args: args{
				accountID:   domain.NewID(uint64(100)),
				operationID: domain.NewID(0),
				amount:      100,
			},
			want:    nil,
			wantErr: domain.NewErrDomain("operation", "'0' is not a valid operation id"),
		},
		{
			name: "repository error",
			fields: fields{
				repo: domain.NewTransactionRepositoryMock(nil, errors.New("repository error")),
			},
			args: args{
				accountID:   domain.NewID(uint64(100)),
				operationID: domain.NewID(4),
				amount:      100,
			},
			want:    nil,
			wantErr: errors.New("repository error"),
		},
		{
			name: "transaction created successfully",
			fields: fields{
				repo: domain.NewTransactionRepositoryMock(domain.NewID(uint64(100)), nil),
			},
			args: args{
				accountID:   domain.NewID(uint64(100)),
				operationID: domain.NewID(4),
				amount:      100,
			},
			want:    transaction.WithID(domain.NewID(uint64(100))),
			wantErr: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCreateTransaction(tt.fields.repo)

			got, err := c.Create(tt.args.accountID, tt.args.operationID, tt.args.amount)
			if (err != nil) && !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			if got.ID().Value() != tt.want.ID().Value() {
				t.Errorf("Invalid Transaction result: ID it must be greater than zero")
				return
			}
		})
	}
}
