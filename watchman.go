package main

import (
	"encoding/json"
	"os"
	"os/exec"
)

func main() {
	cl := NewClient()

	if err := cl.Connect(); err != nil {
		panic(err)
	}

	if err := cl.Close(); err != nil {
		panic(err)
	}
}

func socketLoc() (string, error) {
	if addr := os.Getenv("WATCHMAN_SOCK"); addr != "" {
		return addr, nil
	}

	var loc struct {
		Version  string
		Sockname string
	}

	b, err := exec.Command("watchman", "get-sockname").Output()

	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(b, &loc); err != nil {
		return "", err
	}

	return loc.Sockname, nil
}
