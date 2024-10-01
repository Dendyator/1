package main

import (
	"os/exec"
)

// RunCmd запускает команду + аргументы (cmd) с переменными окружения из env..
func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := exec.Command(cmd[0], cmd[1:]...) //nolint
	command.Env = envToSlice(env)

	if err := command.Start(); err != nil {
		return 1
	}
	if err := command.Wait(); err != nil {
		return 1
	}
	return 0
}

func envToSlice(env Environment) []string {
	var result []string
	envMap := make(map[string]string)

	for k, v := range env {
		if !v.NeedRemove {
			envMap[k] = v.Value
		}
	}

	for k, v := range env {
		if v.NeedRemove {
			delete(envMap, k)
		}
	}

	for k, v := range envMap {
		result = append(result, k+"="+v)
	}

	return result
}
