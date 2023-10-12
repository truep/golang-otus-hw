package main

import (
	"errors"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) < 2 {
		return 1
	}
	for key, value := range env {
		os.Unsetenv(key)
		if !value.NeedRemove {
			os.Setenv(key, value.Value)
		}
	}
	//nolint:gosec
	command := exec.Command(cmd[0], cmd[1:]...)

	command.Env = os.Environ()

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
	}
	exitCode := command.ProcessState.ExitCode()
	return exitCode
}
