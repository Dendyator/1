package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestRunInvalidCmd(t *testing.T) {
	cmd := []string{"invalid_command"}
	env := Environment{}

	returnCode := RunCmd(cmd, env)
	if returnCode == 0 {
		t.Error("Ожидался код завершения >0 для неверной команды")
	}
}

func TestRunCmd(t *testing.T) {
	err := os.Setenv("TEST_VAR", "test_value")
	if err != nil {
		t.Fatalf("Переменная окружения не установлена: %v", err)
	}

	cmd := exec.Command("bash", "-c", "echo $TEST_VAR")

	output, err2 := cmd.CombinedOutput()
	if err2 != nil {
		t.Fatalf("Ошибка при запуске команды: %v", err2)
	}

	expectedOutput := "test_value\n"
	if string(output) != expectedOutput {
		t.Errorf("Ожидаемый результат '%s', полученный результат '%s'", expectedOutput, output)
	}
}
