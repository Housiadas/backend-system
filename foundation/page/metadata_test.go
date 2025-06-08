package page

import (
	"reflect"
	"testing"
)

func Test_calculateMetadata(t *testing.T) {
	type args struct {
		total int
		page  int
		rows  int
	}
	tests := []struct {
		name string
		args args
		want Metadata
	}{
		{
			name: "Metadata empty example",
			args: args{
				total: 0,
				page:  0,
				rows:  0,
			},
			want: Metadata{},
		},
		{
			name: "Metadata example, current page 1",
			args: args{
				total: 150,
				page:  1,
				rows:  8,
			},
			want: Metadata{
				FirstPage:   1,
				CurrentPage: 1,
				LastPage:    19,
				RowsPerPage: 8,
				Total:       150,
			},
		},
		{
			name: "Metadata example, current page 5",
			args: args{
				total: 150,
				page:  5,
				rows:  8,
			},
			want: Metadata{
				FirstPage:   1,
				CurrentPage: 5,
				LastPage:    19,
				RowsPerPage: 8,
				Total:       150,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateMetadata(tt.args.total, tt.args.page, tt.args.rows); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("calculateMetadata() = %v, want %v", got, tt.want)
			}
		})
	}
}
