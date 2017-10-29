// +build windows

package main

import (
	"os"
	osexec "os/exec"
	"os/signal"
	"syscall"
)

func exec(command string, args []string, envv []string) error {
	cmd := osexec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = envv

	sigc := make(chan os.Signal, 1)
	defer close(sigc)
	signal.Notify(sigc, os.Interrupt)

	go func() {
		for sig := range sigc {
			if cmd.Process != nil {
				cmd.Process.Signal(sig)
			}
		}
	}()

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
