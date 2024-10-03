package main

import (
	"os"
	"os/exec"
)

func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := exec.Command(cmd[0], cmd[1:]...) //nolint
	command.Env = envToSlice(env)

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		return 1
	}
	return 0
}

func envToSlice(env Environment) []string {
	var result []string

	for k, v := range env {
		if !v.NeedRemove {
			result = append(result, k+"="+v.Value)
		} else {
			result = append(result, k+"=")
		}
	}

	return result
}
