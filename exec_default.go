// +build !linux

package main

import (
	"fmt"
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

	signals := make([]os.Signal, 32)
	for i := range signals {
		signals[i] = syscall.Signal(i + 1)
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, signals...)

	go func() {
		sig := <-sigc
		if cmd.Process != nil {
			cmd.Process.Signal(sig)
		}
	}()

	var waitStatus syscall.WaitStatus
	if err := cmd.Run(); err != nil {
		if err != nil {
			return err
		}
		if exitError, ok := err.(*osexec.ExitError); ok {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
			os.Exit(waitStatus.ExitStatus())
		}
	}

	return nil
}
