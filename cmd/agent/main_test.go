package main

import (
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)

func TestUpdateStats(t *testing.T) {
	var m Metrics
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
				m:   &m,
				rtm: &rtm,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateStats(tt.args.m, tt.args.rtm)
			assert.NotEmpty(t, m)
			if m.PollCount != tt.want {
				t.Errorf("значение счетчика должно быть %d", tt.want)
			}
		})
	}
}
