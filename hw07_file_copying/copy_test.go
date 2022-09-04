package main

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

const testDataDir = "testdata"

func TestCopy(t *testing.T) {
	testCase := []struct {
		testName            string
		offset              int64
		limit               int64
		expectedContentFile string
	}{
		{testName: "copy full file", offset: 0, limit: 0, expectedContentFile: "out_offset0_limit0.txt"},
		{testName: "copy small amount from begin", offset: 0, limit: 10, expectedContentFile: "out_offset0_limit10.txt"},
		{testName: "copy medium amount from begin", offset: 0, limit: 1000, expectedContentFile: "out_offset0_limit1000.txt"},
		{testName: "copy large amount from begin", offset: 0, limit: 10000, expectedContentFile: "out_offset0_limit10000.txt"},
		{testName: "copy medium amount from middle", offset: 100, limit: 1000, expectedContentFile: "out_offset100_limit1000.txt"},
		{testName: "copy large amount from middle", offset: 6000, limit: 1000, expectedContentFile: "out_offset6000_limit1000.txt"},
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

			fromFilePath := filepath.Join(testDataDir, "input.txt")
			expectedFilePath := filepath.Join(testDataDir, tc.expectedContentFile)
			toFilePath := filepath.Join(tempDir, "to_file.txt")

			copyFile(fromFilePath, toFilePath, tc.offset, tc.limit)

			copiedContent, err := os.ReadFile(toFilePath)
			require.NotErrorIs(t, err, os.ErrNotExist, "target file must exist")

			expectedContent, err := os.ReadFile(expectedFilePath)
			if err != nil {
				t.Fatalf("Can't read expected file '%v'", expectedFilePath)
			}
			require.Equal(t, expectedContent, copiedContent, "contents should be equal")
		})
	}
}

func TestExistedToFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "TestCopy")
	if err != nil {
		t.Fatalf("Error while creating temporary directory")
	}

	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("Error while removing directory %v", tempDir)
		}
	}()

	fromFilePath := filepath.Join(testDataDir, "input.txt")
	expectedFilePath := filepath.Join(testDataDir, "out_offset0_limit10.txt")
	toFilePath := filepath.Join(tempDir, "to_file.txt")

	copyFile(fromFilePath, toFilePath, 0, 0)
	copyFile(fromFilePath, toFilePath, 0, 10)

	copiedContent, err := os.ReadFile(toFilePath)
	require.NoError(t, err, "there shouldn't be an error")

	expectedContent, err := os.ReadFile(expectedFilePath)
	if err != nil {
		t.Fatalf("Can't read expected file '%v'", expectedFilePath)
	}
	require.Equal(t, expectedContent, copiedContent, "contents should be equal")
}

func TestCopyStream(t *testing.T) {
	testCase := []struct {
		testName string
		content  string
		offset   int64
		limit    int64
		expected string
	}{
		{testName: "copy full file", content: "content", expected: "content"},
		{testName: "copy empty file", content: "", expected: ""},
		{testName: "copy from middle to end", content: "content", offset: 3, expected: "tent"},
		{testName: "copy some bytes from begin", content: "content", limit: 3, expected: "con"},
		{testName: "copy some bytes from middle", content: "content", offset: 3, limit: 3, expected: "ten"},
		{testName: "copy large limit", content: "content", limit: 100, expected: "content"},
	}
	for _, tc := range testCase {
		t.Run(tc.testName, func(t *testing.T) {
			fromFile := bytes.NewReader([]byte(tc.content))
			var toFile bytes.Buffer

			err := copyReaderToWriter(fromFile, &toFile, tc.offset, tc.limit)
			require.NoError(t, err, "there shouldn't any errors")
			require.Equal(t, tc.expected, string(toFile.Bytes()), "contents should be equal")
		})
	}
}
