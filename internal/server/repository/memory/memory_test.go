package memory

import (
	"sync"
	"testing"
)

func TestMemStorage_Set(t *testing.T) {
	type fields struct {
		mutex    *sync.RWMutex
		Gauges   map[string]float64
		Counters map[string]int64
	}
	type args struct {
		typ   string
		name  string
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Positive test: counter",
			fields: fields{
				mutex:    &sync.RWMutex{},
				Gauges:   map[string]float64{},
				Counters: map[string]int64{},
			},
			args: args{
				typ:   "counter",
				name:  "test",
				value: "1",
			},
			wantErr: false,
		},
		{
			name: "Positive test: gauge",
			fields: fields{
				mutex:    &sync.RWMutex{},
				Gauges:   map[string]float64{},
				Counters: map[string]int64{},
			},
			args: args{
				typ:   "gauge",
				name:  "test",
				value: "1",
			},
			wantErr: false,
		},
		{
			name: "Negative test: error type",
			fields: fields{
				mutex:    &sync.RWMutex{},
				Gauges:   map[string]float64{},
				Counters: map[string]int64{},
			},
			args: args{
				typ:   "test",
				name:  "test",
				value: "1",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				mutex:    tt.fields.mutex,
				Gauges:   tt.fields.Gauges,
				Counters: tt.fields.Counters,
			}
			if err := m.Set(tt.args.typ, tt.args.name, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
