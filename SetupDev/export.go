package main

import (
	"os"
	"os/exec"
)

type Export struct {
	Key   string
	Value string
}

func (e *Export) Set() error {
	print("Inside\n")
	if past := os.Getenv(e.Key); past == "" {
		print(past, "\n")
		cmd := exec.Command("setx", e.Key, e.Value)
		return cmd.Run()
	} else if past != e.Value {
		print(past, "\n")
		cmd := exec.Command("setx", e.Key, past+";"+e.Value)
		return cmd.Run()
	} else {
		print(past, "\n")
		print(e.Value, "\n")
	}
	return nil
}
