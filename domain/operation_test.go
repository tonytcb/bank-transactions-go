package domain

import (
	"reflect"
	"testing"
)

func TestOperation(t *testing.T) {
	type fields struct {
		id *ID
	}

	tests := []struct {
		name           string
		fields         fields
		wantIsIncoming bool
		wantErr        error
	}{
		// fails
		{
			name: "invalid operation",
			fields: fields{
				id: NewID(10),
			},
			wantErr: NewErrDomain("operation", "10"),
		},
		{
			name: "invalid operation",
			fields: fields{
				id: NewID(100),
			},
			wantErr: NewErrDomain("operation", "100"),
		},

		// success
		{
			name: "valid purchase operation",
			fields: fields{
				id: NewID(1),
			},
			wantIsIncoming: false,
			wantErr:        nil,
		},
		{
			name: "valid purchase operation",
			fields: fields{
				id: NewID(2),
			},
			wantIsIncoming: false,
			wantErr:        nil,
		},
		{
			name: "valid withdraw operation",
			fields: fields{
				id: NewID(3),
			},
			wantIsIncoming: false,
			wantErr:        nil,
		},
		{
			name: "valid payment operation",
			fields: fields{
				id: NewID(4),
			},
			wantIsIncoming: true,
			wantErr:        nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o, err := NewOperation(tt.fields.id)

			if (err != nil) && !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("NewOperation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				if v, ok := err.(*ErrDomain); ok {
					if v.Error() != tt.wantErr.Error() {
						t.Errorf("NewOperation() error.Error() = %v, wantErr %v", err.Error(), tt.wantErr.Error())
						return
					}
				}
			}

			if err != nil {
				return
			}

			if got := o.IsIncoming(); got != tt.wantIsIncoming {
				t.Errorf("IsIncoming() = %v, want %v", got, tt.wantIsIncoming)
			}
		})
	}
}
