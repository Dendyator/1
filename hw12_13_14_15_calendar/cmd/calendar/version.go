package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var (
	release   = "UNKNOWN" //nolint
	buildDate = "UNKNOWN" //nolint
	gitHash   = "UNKNOWN" //nolint
)

func printVersion() { //nolint
	if err := json.NewEncoder(os.Stdout).Encode(struct {
		Release   string
		BuildDate string
		GitHash   string
	}{
		Release:   release,
		BuildDate: buildDate,
		GitHash:   gitHash,
	}); err != nil {
		fmt.Printf("error while decode version info: %v\n", err)
	}
}
