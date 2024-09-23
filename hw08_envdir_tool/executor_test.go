package main

import (
	"os"
	"os/exec"
	"path/filepath"
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
	tempDir, err := os.MkdirTemp("", "envdir_test_run")
	if err != nil {
		t.Fatal(err)
	}
	defer func(path string) {
		err2 := os.RemoveAll(path)
		if err2 != nil {
			t.Fatal(err2)
		}
	}(tempDir)

	if err3 := os.WriteFile(filepath.Join(tempDir, "TEST_VAR"), []byte("test_value"), 0644); err3 != nil {
		t.Fatal(err3)
	}

	env := map[string]EnvValue{
		"TEST_VAR": {"test_value", false},
	}

	cmd := exec.Command("cmd", "/C", "echo %TEST_VAR%")
	cmd.Env = envToSlice(env)

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to run command: %v", err)
	}

	expectedOutput := "test_value\r\n"
	if string(output) != expectedOutput {
		t.Errorf("Ожидаемый результат '%s', получено '%s'", expectedOutput, output)
	}
}
