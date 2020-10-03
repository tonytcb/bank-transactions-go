package domain

import (
	"errors"
	"reflect"
	"testing"
)

func TestNewTransaction(t *testing.T) {
	type args struct {
		accountID   *ID
		operationID *ID
		amount      float64
	}

	tests := []struct {
		name    string
		args    args
		want    *Transaction
		wantErr error
	}{
		// fails
		{
			name: "returns error when the operation 0 is invalid",
			args: args{
				accountID:   NewID(100),
				operationID: NewID(0),
				amount:      201,
			},
			want:    nil,
			wantErr: NewErrDomain("operation", "'0' is not a valid operation id"),
		},
		{
			name: "returns error when the operation 10 is invalid",
			args: args{
				accountID:   NewID(100),
				operationID: NewID(10),
				amount:      201,
			},
			want:    nil,
			wantErr: NewErrDomain("operation", "'10' is not a valid operation id"),
		},

		// successes
		{
			name: "valid transaction with OperationCompraAVista",
			args: args{
				accountID:   NewID(100),
				operationID: NewID(1),
				amount:      201,
			},
			want: &Transaction{
				id: NewID(0),
				account: &Account{
					id: NewID(100),
				},
				operation: OperationCompraAVista,
				amount:    -201,
			},
			wantErr: nil,
		},
		{
			name: "valid transaction with OperationCompraParcelada",
			args: args{
				accountID:   NewID(100),
				operationID: NewID(2),
				amount:      201,
			},
			want: &Transaction{
				id: NewID(0),
				account: &Account{
					id: NewID(100),
				},
				operation: OperationCompraParcelada,
				amount:    -201,
			},
			wantErr: nil,
		},
		{
			name: "valid transaction with OperationSaque",
			args: args{
				accountID:   NewID(100),
				operationID: NewID(3),
				amount:      201,
			},
			want: &Transaction{
				id: NewID(0),
				account: &Account{
					id: NewID(100),
				},
				operation: OperationSaque,
				amount:    -201,
			},
			wantErr: nil,
		},
		{
			name: "valid transaction with OperationPagamento",
			args: args{
				accountID:   NewID(100),
				operationID: NewID(4),
				amount:      201,
			},
			want: &Transaction{
				id: NewID(0),
				account: &Account{
					id: NewID(100),
				},
				operation: OperationPagamento,
				amount:    201,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTransaction(tt.args.accountID, tt.args.operationID, tt.args.amount)

			if (err != nil) && !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("NewTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			if got.Amount() != tt.want.Amount() {
				t.Errorf("NewTransaction() Amount got = %v, want %v", got.Amount(), tt.want.Amount())
			}

			if !reflect.DeepEqual(got.Operation(), tt.want.Operation()) {
				t.Errorf("NewTransaction() Operation() got = %v, want %v", got.Operation(), tt.want.Operation())
			}

			if !reflect.DeepEqual(got.CreatedAt(), tt.want.CreatedAt()) {
				t.Errorf("NewTransaction() CreatedAt() got = %v, want %v", got.CreatedAt(), tt.want.CreatedAt())
			}

			if !reflect.DeepEqual(got.Account().ID(), tt.want.Account().ID()) {
				t.Errorf("NewTransaction() Account.ID got = %v, want %v", got.Account().ID(), tt.want.Account().ID())
			}

			if !reflect.DeepEqual(got.ID(), tt.want.ID()) {
				t.Errorf("NewTransaction() ID got = %v, want %v", got.ID(), tt.want.ID())
			}
		})
	}
}

func TestTransaction_Store(t *testing.T) {
	transaction, _ := NewTransaction(NewID(uint64(1)), NewID(uint64(1)), 100)

	type args struct {
		repo TransactionRepository
	}
	tests := []struct {
		name    string
		args    args
		want    *Transaction
		wantErr error
	}{
		{
			name: "transaction stored successfully",
			args: args{
				repo: NewTransactionRepositoryMock(NewID(uint64(1000)), nil),
			},
			want:    &Transaction{id: NewID(uint64(1000))},
			wantErr: nil,
		},
		{
			name: "error when the repository returns an unknown error",
			args: args{
				repo: NewTransactionRepositoryMock(nil, errors.New("unknown repository error")),
			},
			want:    nil,
			wantErr: errors.New("unknown repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := transaction.Store(tt.args.repo)

			if (err != nil) && !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("Store() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			if got.ID().Value() <= 0 {
				t.Errorf("Invalid Account ID: it must be greater than zero")
				return
			}

			if got.CreatedAt().Year() == 1 { // represents a not defined date
				t.Errorf("Invalid createdAt value: %v", got.CreatedAt().String())
			}
		})
	}
}

func TestTransaction_WithAccount(t1 *testing.T) {
	type args struct {
		a *Account
	}
	tests := []struct {
		name string
		args args
		want *Transaction
	}{
		{
			name: "set account 1",
			args: args{a: &Account{id: NewID(uint64(1))}},
			want: &Transaction{id: NewID(uint64(10))},
		},
		{
			name: "set account 100",
			args: args{a: &Account{id: NewID(uint64(100))}},
			want: &Transaction{id: NewID(uint64(10))},
		},
	}

	for _, tt := range tests {
		transaction := &Transaction{}

		t1.Run(tt.name, func(t1 *testing.T) {
			got := transaction.WithAccount(tt.args.a)

			if transaction == got {
				t1.Error("WithAccount() should return a new Account struct to assure immutability")
			}

			if got.Account().ID().Value() != tt.args.a.ID().Value() {
				t1.Errorf("Account().ID = %v, want %v", got.Account().ID(), tt.args.a.ID())
			}
		})
	}
}
