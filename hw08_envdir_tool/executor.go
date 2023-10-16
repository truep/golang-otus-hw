package main

import (
	"errors"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) < 1 {
		return 1
	}
	var command *exec.Cmd
	name, args := cmd[0], cmd[1:]

	if len(args) == 0 {
		command = exec.Command(name)
	} else {
		command = exec.Command(name, args...)
	}

	env.UdpateEnv()

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
