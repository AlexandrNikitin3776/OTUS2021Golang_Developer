package main

import (
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	environment := make(map[string]EnvValue)

	for _, file := range files {
		if !isValidFileName(file.Name()) {
			continue
		}

		filePath := filepath.Join(dir, file.Name())
		content, err := readFileString(filePath)
		if err != nil {
			return nil, err
		}

		environment[file.Name()] = EnvValue{
			convertContent(content),
			isFileEmpty(content),
		}
	}
	return environment, nil
}

func convertContent(content string) string {
	content = trimSpacesAndTabsRight(content)
	content = replaceHexZeroesToNewLines(content)
	return content
}

func isValidFileName(fileName string) bool {
	if strings.Contains(fileName, "=") {
		return false
	}
	return true
}

func readFileString(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func isFileEmpty(content string) bool {
	if len(content) == 0 {
		return true
	}
	return false
}

func trimSpacesAndTabsRight(content string) string {
	return strings.TrimRight(content, " \t")
}

func replaceHexZeroesToNewLines(content string) string {
	return strings.Replace(content, "\x00", "\n", -1)
}
