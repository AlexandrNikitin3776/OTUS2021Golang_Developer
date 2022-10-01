package main

import (
	"bufio"
	"fmt"
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

func fromStrings(env []string) Environment {
	environment := make(Environment, len(env))
	for _, envString := range env {
		keyValue := strings.Split(envString, `=`)
		environment[keyValue[0]] = EnvValue{Value: keyValue[1], NeedRemove: false}
	}
	return environment
}

func (e Environment) toStrings() []string {
	result := make([]string, 0)
	for envName, envValue := range e {
		result = append(result, fmt.Sprintf("%v=%v", envName, envValue.Value))
	}
	return result
}

func (e Environment) update(target Environment) {
	for envName, envValue := range target {
		if envValue.NeedRemove {
			delete(e, envName)
		} else {
			e[envName] = envValue
		}
	}
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
		content, err := readFirstLine(filePath)
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

func readFirstLine(path string) (string, error) {
	readFile, err := os.Open(path)
	if err != nil {
		return "", err
	}
	scanner := bufio.NewScanner(readFile)
	if ok := scanner.Scan(); ok {
		return scanner.Text(), nil
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", nil
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
