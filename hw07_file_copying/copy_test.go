package main

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"testing"
)

const testDataDir = "testdata"

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
				t.Fatal("Error while creating temporary directory")
			}

			defer func() {
				if err := os.RemoveAll(tempDir); err != nil {
					t.Fatalf("Error while removing directory %v", tempDir)
				}
			}()

			fromFilePath := filepath.Join(testDataDir, "input.txt")
			expectedFilePath := filepath.Join(testDataDir, tc.expectedContentFile)
			toFilePath := filepath.Join(tempDir, "to_file.txt")

			copyErr := Copy(fromFilePath, toFilePath, tc.offset, tc.limit)
			require.NoError(t, copyErr, "there shouldn't be an error")

			copiedContent, err := os.ReadFile(toFilePath)
			require.NotErrorIs(t, err, fs.ErrNotExist, "target file must exist")

			expectedContent, err := os.ReadFile(expectedFilePath)
			if err != nil {
				t.Fatalf("Can't read expected file '%v'", expectedFilePath)
			}
			require.Equal(t, expectedContent, copiedContent, "contents should be equal")
		})
	}
}

func TestExistedToFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "TestExistedToFile")
	if err != nil {
		t.Fatal("Error while creating temporary directory")
	}

	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("Error while removing directory %v", tempDir)
		}
	}()

	fromFilePath := filepath.Join(testDataDir, "input.txt")
	expectedFilePath := filepath.Join(testDataDir, "out_offset0_limit10.txt")
	toFilePath := filepath.Join(tempDir, "to_file.txt")

	copyErr := Copy(fromFilePath, toFilePath, 0, 0)
	require.NoError(t, copyErr, "there shouldn't be an error")

	copyErr = Copy(fromFilePath, toFilePath, 0, 10)
	require.NoError(t, copyErr, "there shouldn't be an error")

	copiedContent, err := os.ReadFile(toFilePath)
	require.NotErrorIs(t, err, fs.ErrNotExist, "target file must exist")

	expectedContent, err := os.ReadFile(expectedFilePath)
	if err != nil {
		t.Fatalf("Can't read expected file '%v'", expectedFilePath)
	}
	require.Equal(t, expectedContent, copiedContent, "contents should be equal")
}

func TestCopyDevRandom(t *testing.T) {
	var bytesToRead int64 = 1000
	tempDir, err := os.MkdirTemp("", "TestCopyDevRandom")
	if err != nil {
		t.Fatal("Error while creating temporary directory")
	}

	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("Error while removing directory %v", tempDir)
		}
	}()

	fromFilePath := "/dev/urandom"
	toFilePath := filepath.Join(tempDir, "to_file.txt")

	copyErr := Copy(fromFilePath, toFilePath, 0, bytesToRead)
	require.NoError(t, copyErr, "there shouldn't be an error")

	copiedContent, err := os.ReadFile(toFilePath)
	require.NotErrorIs(t, err, fs.ErrNotExist, "target file must exist")
	require.Equal(t, int64(len(copiedContent)), bytesToRead)
}

func TestOffsetExceedsFileSizeFail(t *testing.T) {
	content := "some content"
	tempDir, err := os.MkdirTemp("", "TestOffsetExceedsFileSizeFail")
	if err != nil {
		t.Fatal("Error while creating temporary directory")
	}

	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("Error while removing directory %v", tempDir)
		}
	}()

	fromFilePath := filepath.Join(tempDir, "from_.txt")

	if err = os.WriteFile(fromFilePath, []byte(content), os.ModePerm); err != nil {
		t.Fatalf("Error while writing temp file %v", fromFilePath)
	}

	toFilePath := filepath.Join(tempDir, "to_file.txt")

	copyErr := Copy(fromFilePath, toFilePath, math.MaxInt64, 0)
	require.ErrorIs(t, copyErr, ErrOffsetExceedsFileSize, "offset larger than file size is invalid")

	_, err = os.ReadFile(toFilePath)
	require.ErrorIs(t, err, fs.ErrNotExist, "target file should not exist")
}

func TestFromFileNotExist(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "TestFromFileNotExist")
	if err != nil {
		t.Fatal("Error while creating temporary directory")
	}

	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("Error while removing directory %v", tempDir)
		}
	}()

	fromFilePath := filepath.Join(testDataDir, "not_existed_file.txt")
	toFilePath := filepath.Join(tempDir, "to_file.txt")

	copyErr := Copy(fromFilePath, toFilePath, math.MaxInt64, 0)
	require.ErrorIs(t, copyErr, fs.ErrNotExist, "from file doesn't exist")

	_, err = os.ReadFile(toFilePath)
	require.ErrorIs(t, err, fs.ErrNotExist, "target file should not exist")
}
