package main

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestParseCLIArgsOK(t *testing.T) {
	expectedDir := "dir"
	expectedCmd := []string{"command", "arg1", "arg2"}
	os.Args = append([]string{"test", expectedDir}, expectedCmd...)
	dir, cmd, err := ParseCLIArgs()
	require.NoError(t, err, "shouldn't be any error")
	require.Equal(t, expectedDir, dir)
	require.Equal(t, expectedCmd, cmd)
}

func TestParseCLIArgsFail(t *testing.T) {
	os.Args = []string{"test", "dir"}
	_, _, err := ParseCLIArgs()
	require.Error(t, err, "should be an error")
}
