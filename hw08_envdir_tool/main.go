package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go-envdir /path/to/env/dir command arg1 arg2 ...")
		os.Exit(1)
	}

	envDir := os.Args[1]
	cmdArgs := os.Args[2:]
	env, err := ReadDir(envDir)
	if err != nil {
		fmt.Printf("Error setting environment variables: %v\n", err)
		os.Exit(1)
	}

	returnCode := RunCmd(cmdArgs, env)
	os.Exit(returnCode)
}
