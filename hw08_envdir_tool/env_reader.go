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

		lines := strings.SplitN(string(content), "\n", 2)
		firstLine := lines[0]
		cleanedValue := strings.TrimRight(firstLine, " \t")

		cleanedValue = strings.ReplaceAll(cleanedValue, "\x00", "")

		env[file.Name()] = EnvValue{Value: cleanedValue, NeedRemove: len(cleanedValue) == 0}
	}

	return env, nil
}
