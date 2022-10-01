package main

import (
	"log"
	"os"
)

func main() {
	dir, cmd, err := ParseCLIArgs()
	if err != nil {
		log.Fatalf("Can't parse CLI arguments: %v", err)
	}
	env, err := ReadDir(dir)
	if err != nil {
		log.Fatalf("Can't parse env directory: %v", err)
	}
	returnCode := RunCmd(cmd, env)
	os.Exit(returnCode)
}
