package domain

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestNewAccount(t *testing.T) {
	type args struct {
		documentNumber DocumentNumber
	}

	tests := []struct {
		name    string
		args    args
		want    *Account
		wantErr error
	}{
		{
			name:    "invalid document number",
			args:    args{documentNumber: "00000000000"},
			want:    nil,
			wantErr: NewErrDomain("document.number", "'00000000000' is not a valid document number"),
		},
		{
			name: "valid CPF document",
			args: args{documentNumber: "00000000191"},
			want: &Account{
				id:       NewID(uint64(0)),
				document: &Document{number: DocumentNumber("00000000191")},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAccount(tt.args.documentNumber)

			if (err != nil) && !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("NewAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			if !reflect.DeepEqual(got.ID(), tt.want.ID()) {
				t.Errorf("NewAccount().ID() got = %v, want %v", got.ID(), tt.want.ID())
				return
			}

			if !reflect.DeepEqual(got.Document(), tt.want.Document()) {
				t.Errorf("NewAccount().Document() got = %v, want %v", got.Document(), tt.want.Document())
			}
		})
	}
}

func TestAccount_Store(t *testing.T) {
	baseAccount := &Account{
		id:       NewID(0),
		document: &Document{number: "00000000191"},
	}

	type args struct {
		repo AccountRepositoryWriter
	}

	tests := []struct {
		name    string
		args    args
		want    *Account
		wantErr error
	}{
		{
			name: "store successful",
			args: args{
				repo: NewAccountRepositoryMock(NewID(uint64(1)), nil, nil),
			},
			want: &Account{
				id: NewID(uint64(1)),
			},
			wantErr: nil,
		},
		{
			name: "store error",
			args: args{
				repo: NewAccountRepositoryMock(nil, nil, errors.New("unknown repository error")),
			},
			want:    nil,
			wantErr: errors.New("unknown repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := baseAccount.Store(tt.args.repo)

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

func TestAccount_WithID(t *testing.T) {
	type args struct {
		id *ID
	}
	tests := []struct {
		name string
		args args
		want *ID
	}{
		{
			name: "set id 100 successfully",
			args: args{
				id: NewID(100),
			},
			want: NewID(100),
		},
		{
			name: "set id 1 successfully",
			args: args{
				id: NewID(1),
			},
			want: NewID(1),
		},
	}

	for _, tt := range tests {
		a := &Account{
			id:       NewID(1),
			document: &Document{number: "00000000191"},
		}

		t.Run(tt.name, func(t *testing.T) {
			got := a.WithID(tt.args.id)

			if a == got {
				t.Error("WithID() should return a new Account struct to assure immutability")
			}

			if got.ID().Value() != tt.want.Value() {
				t.Errorf("ID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccount_WithCreateAt(t *testing.T) {
	var (
		yesterday = time.Now().Add(-24 * time.Hour)
		tomorrow  = time.Now().Add(24 * time.Hour)
	)

	type args struct {
		createdAt time.Time
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "set yesterday as createdAt",
			args: args{
				createdAt: yesterday,
			},
			want: yesterday,
		},
		{
			name: "set tomorrow as createdAt",
			args: args{
				createdAt: tomorrow,
			},
			want: tomorrow,
		},
	}

	for _, tt := range tests {
		a := &Account{
			createdAt: time.Now(),
		}

		t.Run(tt.name, func(t *testing.T) {
			got := a.WithCreateAt(tt.args.createdAt)

			if a == got {
				t.Error("WithCreateAt() should return a new Account struct to assure immutability")
			}

			if got.CreatedAt().String() != tt.want.String() {
				t.Errorf("CreatedAt() = %v, want %v", got, tt.want)
			}
		})
	}
}
