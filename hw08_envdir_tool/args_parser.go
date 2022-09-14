package main

import (
	"fmt"
	"os"
)

type ErrInvalidUsage struct {
	app string
}

func (e *ErrInvalidUsage) Error() string {
	return fmt.Sprintf("usage of %v is : '%v dir command [args]'", e.app, e.app)
}

type CLIArgs struct {
	dir     string
	command string
	args    []string
}

func ParseCLIArgs() (*CLIArgs, error) {
	if len(os.Args) < 3 {
		return nil, &ErrInvalidUsage{os.Args[0]}
	}
	args := make([]string, 0)
	if len(os.Args) > 3 {
		args = os.Args[3:len(os.Args)]
	}
	return &CLIArgs{os.Args[1], os.Args[2], args}, nil
}
