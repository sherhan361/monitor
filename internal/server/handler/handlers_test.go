package handler

import (
	"net/http"
	"testing"
)

func Test_checkParams(t *testing.T) {
	type args struct {
		typ   string
		name  string
		value string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Positive test: counter",
			args: args{
				typ:   "counter",
				name:  "test",
				value: "1",
			},
			want: http.StatusOK,
		},
		{
			name: "Positive test: gauge",
			args: args{
				typ:   "gauge",
				name:  "test",
				value: "1.1",
			},
			want: http.StatusOK,
		},
		{
			name: "Negative test: error type",
			args: args{
				typ:   "test",
				name:  "test",
				value: "1",
			},
			want: http.StatusNotImplemented,
		},
		{
			name: "Negative test: empty params",
			args: args{
				typ:   "",
				name:  "",
				value: "",
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Negative test: none counter",
			args: args{
				typ:   "counter",
				name:  "test",
				value: "none",
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Negative test: none gauge",
			args: args{
				typ:   "gauge",
				name:  "test",
				value: "none",
			},
			want: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkParams(tt.args.typ, tt.args.name, tt.args.value); got != tt.want {
				t.Errorf("checkParams() = %v, want %v", got, tt.want)
			}
		})
	}
}
