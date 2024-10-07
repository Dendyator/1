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

		fileInfo, err := file.Info()
		if err != nil || !fileInfo.Mode().IsRegular() {
			continue
		}

		content, err2 := os.ReadFile(filepath.Join(dir, file.Name()))
		if err2 != nil {
			return nil, err2
		}

		value := strings.TrimRight(string(content), "\n \t")
		cleanedValue := strings.ReplaceAll(value, "\x00", "\n")
		env[file.Name()] = EnvValue{Value: cleanedValue, NeedRemove: len(cleanedValue) == 0}
	}

	return env, nil
}
