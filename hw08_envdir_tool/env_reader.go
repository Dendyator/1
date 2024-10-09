package main

import (
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

type EnvValue struct {
	Value      string
	NeedRemove bool
}

func ReadDir(dir string) (Environment, error) {
	env := make(Environment)

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() || strings.Contains(file.Name(), "=") {
			continue
		}

		filePath := filepath.Join(dir, file.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		value := strings.ReplaceAll(string(content), "\x00", "")
		value = strings.TrimRight(value, " \t")
		
		env[file.Name()] = EnvValue{Value: value, NeedRemove: len(value) == 0}
	}

	return env, nil
}
