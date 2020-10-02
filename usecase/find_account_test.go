package usecase

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/tonytcb/bank-transactions-go/infra/repository"

	"github.com/tonytcb/bank-transactions-go/domain"
)

func TestFindAccount_Find(t *testing.T) {
	accountOK, _ := domain.NewAccount("00000000191")
	accountOK = accountOK.WithID(domain.NewID(uint64(100))).WithCreateAt(time.Now())

	type fields struct {
		repo domain.AccountRepository
	}
	type args struct {
		id *domain.ID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Account
		wantErr error
	}{
		{
			name: "account not found error",
			fields: fields{
				repo: domain.NewAccountMock(nil, nil, repository.NewErrRegisterNotFound("account", "100")),
			},
			args: args{
				id: domain.NewID(uint64(100)),
			},
			want:    nil,
			wantErr: repository.NewErrRegisterNotFound("account", "100"),
		},
		{
			name: "unknown repository error",
			fields: fields{
				repo: domain.NewAccountMock(nil, nil, errors.New("some repository error")),
			},
			args: args{
				id: domain.NewID(uint64(100)),
			},
			want:    nil,
			wantErr: errors.New("some repository error"),
		},
		{
			name: "account found successfully",
			fields: fields{
				repo: domain.NewAccountMock(nil, accountOK, nil),
			},
			args: args{
				id: domain.NewID(uint64(100)),
			},
			want:    accountOK,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFindAccount(tt.fields.repo)

			got, err := f.Find(tt.args.id)
			if (err != nil) && !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Find() got = %v, want %v", got, tt.want)
			}
		})
	}
}
