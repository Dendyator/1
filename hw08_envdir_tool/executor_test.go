package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestRunCmd(t *testing.T) {
	err := os.Setenv("TEST_VAR", "test_value")
	if err != nil {
		t.Fatalf("Переменная окружения не установлена: %v", err)
	}

	cmd := exec.Command("cmd", "/C", "echo %TEST_VAR%")

	output, err2 := cmd.CombinedOutput()
	if err2 != nil {
		t.Fatalf("Ошибка при запуске команды: %v", err2)
	}

	expectedOutput := "test_value\r\n"
	if string(output) != expectedOutput {
		t.Errorf("Ожидаемый результат '%s', полученный результат '%s'", expectedOutput, output)
	}
}
