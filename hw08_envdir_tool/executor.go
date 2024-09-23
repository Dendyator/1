package main

import (
	"os/exec"
)

// RunCmd запускает команду + аргументы (cmd) с переменными окружения из env..
func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := exec.Command(cmd[0], cmd[1:]...) //nolint

	command.Env = envToSlice(env)

	err := command.Start()
	if err != nil {
		return 1
	}

	err = command.Wait()
	if err != nil {
		return 1
	}

	return 0
}

func envToSlice(env Environment) []string {
	var result []string
	for k, v := range env {
		if !v.NeedRemove {
			result = append(result, k+"="+v.Value)
		}
	}
	return result
}
