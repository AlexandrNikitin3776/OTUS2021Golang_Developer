package main

import (
	"fmt"
	"os"
)

type ErrInvalidUsage struct {
	appName string
}

func (e *ErrInvalidUsage) Error() string {
	return fmt.Sprintf("usage of %v is : '%v dir command [args]'", e.appName, e.appName)
}

func ParseCLIArgs() (string, []string, error) {
	if len(os.Args) < 3 {
		return "", nil, &ErrInvalidUsage{os.Args[0]}
	}
	args := os.Args[2:len(os.Args)]
	return os.Args[1], args, nil
}
