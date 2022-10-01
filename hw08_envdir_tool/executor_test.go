package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRunCmd(t *testing.T) {
	type args struct {
		cmd []string
		env Environment
	}
	tests := []struct {
		name           string
		args           args
		wantReturnCode int
	}{
		{
			"echo env",
			args{
				[]string{"printenv", "NAME"},
				Environment{"NAME": {"42", false}}},
			0,
		},
		{
			"return 56",
			args{
				[]string{"sh", "-c", "exit 56"},
				Environment{}},
			56,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotReturnCode := RunCmd(tt.args.cmd, tt.args.env)
			require.Equal(t, tt.wantReturnCode, gotReturnCode, "values must be equal")

		})
	}
}

func TestSetEnvs(t *testing.T) {
	type args struct {
		commandEnv        []string
		targetEnvironment Environment
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"do nothing",
			args{
				[]string{"a=15"},
				make(Environment, 0),
			},
			[]string{"a=15"},
		},
		{
			"add new",
			args{
				[]string{"a=15"},
				Environment{"b": {"12", false}},
			},
			[]string{"a=15", "b=12"},
		},
		{
			"replace",
			args{
				[]string{"a=15"},
				Environment{"a": {"12", false}},
			},
			[]string{"a=12"},
		},
		{
			"remove",
			args{
				[]string{"a=15"},
				Environment{"a": {"", true}},
			},
			[]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := setEnvs(tt.args.commandEnv, tt.args.targetEnvironment)
			require.ElementsMatch(t, tt.want, got, "values must be equal")
		})
	}
}
