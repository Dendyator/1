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

		content, err := os.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}

		value := strings.TrimRight(string(content), "\n\t ")
		cleanedValue := strings.ReplaceAll(value, "\x00", "")
		env[file.Name()] = EnvValue{Value: cleanedValue, NeedRemove: len(cleanedValue) == 0}
	}

	return env, nil
}
