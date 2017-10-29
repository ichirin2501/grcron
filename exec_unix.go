// +build !windows

package main

import (
	osexec "os/exec"
	"syscall"
)

func exec(command string, args []string, envv []string) error {
	binary, err := osexec.LookPath(command)
	if err != nil {
		return err
	}

	argv := make([]string, 0, 1+len(args))
	argv = append(argv, command)
	argv = append(argv, args...)

	return syscall.Exec(binary, argv, envv)
}
