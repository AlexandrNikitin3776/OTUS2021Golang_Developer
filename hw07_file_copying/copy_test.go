package main

import (
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func cleanUp(t *testing.T, dirName string) {
	if err := os.RemoveAll(dirName); err != nil {
		t.Fatalf("Error while removing directory %v", dirName)
	}
}

func TestCopy(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "TestCopy")
	if err != nil {
		t.Fatalf("Error while creating temporary directory")
	}
	defer cleanUp(t, tempDir)

	fromFilePath := filepath.Join(tempDir, "from_file.txt")
	toFilePath := filepath.Join(tempDir, "to_file.txt")

	t.Run("copy full file", func(t *testing.T) {
		content := []byte("content")

		err := os.WriteFile(fromFilePath, content, 777)
		if err != nil {
			t.Fatalf("Error while writing %v: %v\n", fromFilePath, err)
		}

		dd(fromFilePath, toFilePath, int64(0), int64(0))

		copiedContent, err := os.ReadFile(toFilePath)
		require.NotErrorIs(t, err, os.ErrNotExist, "target file must exist")
		require.Equal(t, content, copiedContent, "contents should be equal")
	})
}
