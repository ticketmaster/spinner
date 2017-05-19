package main

import (
	"testing"

	"golang.org/x/sys/windows/svc"
)

func Test_Runing_testState(t *testing.T) {
	type args struct {
		s svc.State
	}
	tests := []struct {
		name string
		args args
	}{
		{"Running", args{svc.Running}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := testState(tt.args.s)
			if err != nil {
				t.Error(
					"For", tt.name,
					"expected no Error",
					"got", err.Error(),
				)
			}
		})
	}
}

func Test_Terminating_testState(t *testing.T) {
	type args struct {
		s svc.State
	}
	tests := []struct {
		name string
		args args
	}{
		{"Stopped", args{svc.Stopped}},
		{"StartPending", args{svc.StartPending}},
		{"StopPending", args{svc.StopPending}},
		{"ContinuePending", args{svc.ContinuePending}},
		{"Paused", args{svc.Paused}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := testState(tt.args.s)
			if err == nil {
				t.Error(
					"For", tt.name,
					"expected an Error [", tt.args.s,
					"] got nil",
				)
			}

		})
	}
}
