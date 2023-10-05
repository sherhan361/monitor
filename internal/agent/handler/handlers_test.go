package handler

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestUpdateMetrics(t *testing.T) {
	type values struct {
		*sync.RWMutex
		m *Metrics
	}

	tests := []struct {
		name string
		args values
		want int64
	}{
		{
			name: "Test 1",
			args: values{
				m: &Metrics{
					RWMutex:  &sync.RWMutex{},
					Gauges:   map[string]float64{},
					Counters: map[string]int64{},
				},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateMetrics(tt.args.m)
			assert.NotEmpty(t, tt.args.m)
			if tt.args.m.Counters["PollCount"] != tt.want {
				t.Errorf("значение счетчика должно быть %d", tt.want)
			}
		})
	}
}
