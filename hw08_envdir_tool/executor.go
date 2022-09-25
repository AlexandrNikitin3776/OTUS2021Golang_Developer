package main

import (
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) int {
	command := exec.Command(cmd[0], cmd[1:]...)
	command.Env = setEnvs(os.Environ(), env)

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err := command.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode()
		}
		log.Fatalf("comand raises error: %v", err)
	}
	return 0
}

func setEnvs(commandEnv []string, targetEnvironment Environment) []string {
	commandEnvironment := fromStrings(commandEnv)
	commandEnvironment.update(targetEnvironment)
	return commandEnvironment.toStrings()
}
