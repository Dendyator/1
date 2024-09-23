package main

import (
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue помогает различать пустые файлы и файлы с первой пустой строкой.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir считывает указанный каталог и возвращает отображение переменных env
// Переменные представлены в виде файлов, где filename - это имя переменной, а первая строка файла - это значение.
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

		content, err2 := os.ReadFile(filepath.Join(dir, file.Name()))
		if err2 != nil {
			return nil, err2
		}

		value := strings.TrimRight(string(content), "\n \t")
		if len(value) == 0 {
			env[file.Name()] = EnvValue{NeedRemove: true}
		} else {
			env[file.Name()] = EnvValue{Value: strings.ReplaceAll(value, "\x00", "\n"), NeedRemove: false}
		}
	}

	return env, nil
}
