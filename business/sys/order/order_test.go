package order

import (
	"reflect"
	"testing"
)

func Test_Order_Parse(t *testing.T) {
	fieldMappings := map[string]string{
		"user_id": "user_id",
		"name":    "name",
		"email":   "email",
	}
	type args struct {
		orderBy      string
		defaultOrder By
	}
	tests := []struct {
		name    string
		args    args
		want    By
		wantErr bool
	}{
		{
			name: "Default order",
			args: args{
				orderBy: "",
				defaultOrder: By{
					Field:     "user_id",
					Direction: ASC,
				},
			},
			want: By{
				Field:     "user_id",
				Direction: ASC,
			},
			wantErr: false,
		},
		{
			name: "Valid request order by DESC",
			args: args{
				orderBy: "name,DESC",
				defaultOrder: By{
					Field:     "name",
					Direction: DESC,
				},
			},
			want: By{
				Field:     "name",
				Direction: DESC,
			},
			wantErr: false,
		},
		{
			name: "Invalid request order by ASC",
			args: args{
				orderBy: "unknown,DESC",
				defaultOrder: By{
					Field:     "user_id",
					Direction: ASC,
				},
			},
			want:    By{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(fieldMappings, tt.args.orderBy, tt.args.defaultOrder)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}
