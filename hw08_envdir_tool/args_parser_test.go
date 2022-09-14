package main

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestParseCLIArgsOK(t *testing.T) {
	expected := CLIArgs{"dir", "command", []string{"arg1", "arg2"}}
	os.Args = []string{"test", "dir", "command", "arg1", "arg2"}
	result, err := ParseCLIArgs()
	require.NoError(t, err, "shouldn't be any error")
	require.Equal(t, expected, *result)
}

func TestParseCLIArgsFail(t *testing.T) {
	os.Args = []string{"test", "dir"}
	_, err := ParseCLIArgs()
	require.Error(t, err, "should be an error")
}
