package main

import (
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestCopy(t *testing.T) {
	testCase := []struct {
		testName string
		content  string
		offset   int64
		limit    int64
		expected string
	}{
		{testName: "copy full file", content: "content", expected: "content"},
		{testName: "copy empty file", content: "", expected: ""},
		{testName: "copy from file middle", content: "content", offset: 3, expected: "tent"},
	}
	for _, tc := range testCase {
		t.Run(tc.testName, func(t *testing.T) {
			tempDir, err := os.MkdirTemp("", "TestCopy")
			if err != nil {
				t.Fatalf("Error while creating temporary directory")
			}

			defer func() {
				if err := os.RemoveAll(tempDir); err != nil {
					t.Fatalf("Error while removing directory %v", tempDir)
				}
			}()

			fromFilePath := filepath.Join(tempDir, "from_file.txt")
			toFilePath := filepath.Join(tempDir, "to_file.txt")

			content := []byte(tc.content)
			if err := os.WriteFile(fromFilePath, content, 777); err != nil {
				t.Fatalf("Error while writing %v: %v\n", fromFilePath, err)
			}

			dd(fromFilePath, toFilePath, tc.offset, tc.limit)

			copiedContent, err := os.ReadFile(toFilePath)
			require.NotErrorIs(t, err, os.ErrNotExist, "target file must exist")
			require.Equal(t, tc.expected, string(copiedContent), "contents should be equal")
		})
	}
}
