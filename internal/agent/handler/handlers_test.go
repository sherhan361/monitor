package handler

import (
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)

func TestUpdateMetrics(t *testing.T) {
	var rtm runtime.MemStats
	type values struct {
		m   *Metrics
		rtm *runtime.MemStats
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
					Gauges:   map[string]float64{},
					Counters: map[string]int64{},
				},
				rtm: &rtm,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateMetrics(tt.args.m, tt.args.rtm)
			assert.NotEmpty(t, tt.args.m)
			if tt.args.m.Counters["PollCount"] != tt.want {
				t.Errorf("значение счетчика должно быть %d", tt.want)
			}
		})
	}
}
