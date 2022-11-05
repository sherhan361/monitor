package main

import "testing"

func TestReportSender(t *testing.T) {
	type args struct {
		m              *Metrics
		reportInterval int
	}
	tests := []struct {
		name string
		args args
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ReportSender(tt.args.m, tt.args.reportInterval)
		})
	}
}
