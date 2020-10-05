package domain

import (
	"reflect"
	"testing"
)

func TestNewDocument(t *testing.T) {
	type args struct {
		number DocumentNumber
	}

	tests := []struct {
		name    string
		args    args
		want    *Document
		wantErr error
	}{
		{
			name:    "empty document number",
			args:    args{number: ""},
			want:    nil,
			wantErr: NewErrDomain("document.number", "'' is not a valid document number"),
		},
		{
			name:    "document number with one character",
			args:    args{number: "1"},
			want:    nil,
			wantErr: NewErrDomain("document.number", "'1' is not a valid document number"),
		},
		{
			name:    "document number with twenty characters",
			args:    args{number: "12312312312312312312"},
			want:    nil,
			wantErr: NewErrDomain("document.number", "'12312312312312312312' is not a valid document number"),
		},
		{
			name:    "document number with alpha characters",
			args:    args{number: "abcdefghijk"},
			want:    nil,
			wantErr: NewErrDomain("document.number", "'abcdefghijk' is not a valid document number"),
		},
		{
			name:    "invalid document",
			args:    args{number: "11111111111"},
			want:    nil,
			wantErr: NewErrDomain("document.number", "'11111111111' is not a valid document number"),
		},
		{
			name:    "invalid document",
			args:    args{number: "00000000190"},
			want:    nil,
			wantErr: NewErrDomain("document.number", "'00000000190' is not a valid document number"),
		},
		{
			name:    "invalid document",
			args:    args{number: "000000001911"},
			want:    nil,
			wantErr: NewErrDomain("document.number", "'000000001911' is not a valid document number"),
		},
		{
			name:    "valid document",
			args:    args{number: "00000000191"},
			want:    &Document{number: "00000000191"},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDocument(tt.args.number)
			if (err != nil) && !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("NewDocument() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			if got.Number().String() != tt.want.Number().String() {
				t.Errorf("NewDocument() got Number = %v, want %v", got.Number().String(), tt.want.Number().String())
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDocument() got = %v, want %v", got, tt.want)
			}
		})
	}
}
