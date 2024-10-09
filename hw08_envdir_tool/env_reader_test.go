package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadDir(t *testing.T) {
	tempDir := t.TempDir()

	// Создаём файл FOO с пробелами перед значением
	writeFile(t, filepath.Join(tempDir, "FOO"), []byte("    foo\nwith new line"))
	writeFile(t, filepath.Join(tempDir, "BAR"), []byte("value\nignored line"))

	env, err := ReadDir(tempDir)
	if err != nil {
		t.Fatalf("ReadDir() error: %v", err)
	}

	expected := Environment{
		"FOO": {Value: "foo\nwith new line", NeedRemove: false},
	}

	for k, v := range expected {
		if got, ok := env[k]; !ok || got != v {
			t.Errorf("Ожидаемый результат %v для %s, получено %v", v, k, got)
		}
	}
}

func writeFile(t *testing.T, path string, content []byte) {
	t.Helper()
	err := os.WriteFile(path, content, 0644)
	if err != nil {
		t.Fatalf("Файл не записан: %v", err)
	}
}
