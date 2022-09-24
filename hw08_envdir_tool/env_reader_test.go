package main

import (
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

type EnvFile struct {
	name    string
	content string
}

func TestReadDir(t *testing.T) {
	testCases := []struct {
		name     string
		envFiles []EnvFile
		expected Environment
	}{
		{
			"empty",
			nil,
			make(Environment)},
		{
			"env to replace",
			[]EnvFile{{"FOO", "19"}},
			Environment{"FOO": EnvValue{"19", false}},
		},
		{
			"env to remove",
			[]EnvFile{{"BAR", ""}},
			Environment{"BAR": EnvValue{"", true}},
		},
		{
			"with spaces in the end",
			[]EnvFile{{"FOO", "19     "}},
			Environment{"FOO": EnvValue{"19", false}},
		},
		{
			"with tabs in the end",
			[]EnvFile{{"FOO", "19\t\t\t\t\t"}},
			Environment{"FOO": EnvValue{"19", false}},
		},
		{
			"replace terminal zeroes to \n",
			[]EnvFile{{"FOO", "19\x0042"}},
			Environment{"FOO": EnvValue{"19\n42", false}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tempDir, err := os.MkdirTemp(os.TempDir(), tc.name)
			if err != nil {
				t.Fatal("Error while creating temporary directory")
			}

			defer func() {
				if err := os.RemoveAll(tempDir); err != nil {
					t.Fatalf("Error while removing directory %v: %v", tempDir, err)
				}
			}()

			for _, envFile := range tc.envFiles {
				filePath := filepath.Join(tempDir, envFile.name)
				if err = os.WriteFile(filePath, []byte(envFile.content), os.ModePerm); err != nil {
					t.Fatalf("Can't create temporary file: %v", err)
				}
			}

			result, err := ReadDir(tempDir)
			t.Log(result, err)
			require.NoError(t, err, "shouldn't be any error")
			require.Equal(t, tc.expected, result)
		})
	}

}

func TestIsValidFileName(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		want     bool
	}{
		{"valid filename", "FOO", true},
		{"digit filename", "3", true},
		{"invalid filename", "FO0=123", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidFileName(tt.fileName); got != tt.want {
				t.Errorf("isValidFileName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsFileEmpty(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{"empty", "", true},
		{"not empty", "not empty content", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isFileEmpty(tt.content); got != tt.want {
				t.Errorf("isFileEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}
