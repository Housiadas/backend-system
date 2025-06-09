package money

import (
	"reflect"
	"testing"
)

func TestMoney_Equal(t *testing.T) {
	type fields struct {
		value float64
	}

	type args struct {
		m2 Money
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Equal",
			fields: fields{
				value: 100,
			},
			args: args{
				m2: Money{
					value: 100,
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Money{
				value: tt.fields.value,
			}
			if got := m.Equal(tt.args.m2); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoney_String(t *testing.T) {
	type fields struct {
		value float64
	}

	type args struct {
		m2 Money
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "String",
			fields: fields{
				value: 100.55,
			},
			args: args{
				m2: Money{
					value: 100.55,
				},
			},
			want: "100.55",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Money{
				value: tt.fields.value,
			}
			if got := m.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoney_Value(t *testing.T) {
	type fields struct {
		value float64
	}

	type args struct {
		m2 Money
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{
			name: "Value",
			fields: fields{
				value: 100.55,
			},
			args: args{
				m2: Money{
					value: 100.55,
				},
			},
			want: 100.55,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Money{
				value: tt.fields.value,
			}
			if got := m.Value(); got != tt.want {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMustParse(t *testing.T) {
	type args struct {
		m Money
	}

	tests := []struct {
		name string
		args args
		want Money
	}{
		{
			name: "MustParse",
			args: args{
				m: Money{
					value: 100_000,
				},
			},
			want: Money{
				value: 100_000,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MustParse(tt.args.m.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MustParse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	type args struct {
		m Money
	}

	tests := []struct {
		name    string
		args    args
		want    Money
		wantErr bool
	}{
		{
			name: "Parse",
			args: args{
				m: Money{
					value: 100_000,
				},
			},
			want: Money{
				value: 100_000,
			},
			wantErr: false,
		},
		{
			name: "NotParse",
			args: args{
				m: Money{
					value: 10_000_000,
				},
			},
			want:    Money{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.m.value)
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
