/* grcron
 *
 * Copyright (c) 2017 Motoaki Nishikawa
 * Distributed under MIT license, see LICENSE file.
 */

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

const version string = "0.0.4"

type grcron struct {
	StateFile    string
	DefaultState string
	CurrentState string
}

func newGrcron(defaultState string, stateFile string) (*grcron, error) {
	if !(defaultState == "active" || defaultState == "passive") {
		return nil, fmt.Errorf("The Value of DefaultState:%s is incorrect", defaultState)
	}

	f, err := os.Open(stateFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var curr string
	sc := bufio.NewScanner(f)
	sc.Scan()
	st := sc.Text()
	switch st {
	case "active", "passive":
		curr = st
	default:
		fmt.Fprintf(os.Stderr, "corrupted state file('%s') (content='%s'), staying at gr.DefaultState('%s')\n", stateFile, st, defaultState)
		curr = defaultState
	}

	return &grcron{
		DefaultState: defaultState,
		StateFile:    stateFile,
		CurrentState: curr,
	}, nil
}

var testKeepalivedActive func() (bool, error)

func (gr grcron) keepalivedActive() (bool, error) {
	if testKeepalivedActive != nil {
		return testKeepalivedActive()
	}
	cmd := exec.Command("sh", "-c", "ps cax | grep -q keepalived")
	err := cmd.Run()
	// 異常終了はkeepalivedプロセスがいないとみなす
	if _, ok := err.(*exec.ExitError); ok {
		return false, fmt.Errorf("keepalived is probably down")
	}
	return true, nil
}

func (gr grcron) canRun() (bool, error) {
	ka, err := gr.keepalivedActive()
	if err != nil {
		return false, err
	}
	return gr.CurrentState == "active" && ka, nil
}

func main() {
	var (
		showVersion  bool
		dryRun       bool
		stateFile    string
		defaultState string
	)

	flag.StringVar(&stateFile, "f", "/var/lib/grcron/state", "grcron state file.")
	flag.StringVar(&defaultState, "s", "passive", "grcron default state.")
	flag.BoolVar(&showVersion, "version", false, "show version number.")
	flag.BoolVar(&showVersion, "v", false, "show version number.")
	flag.BoolVar(&dryRun, "dryrun", false, "dry-run.")
	flag.BoolVar(&dryRun, "n", false, "dry-run.")
	flag.Parse()
	args := flag.Args()

	if showVersion {
		fmt.Printf("grcron %s, %s built for %s/%s\n", version, runtime.Version(), runtime.GOOS, runtime.GOARCH)
		return
	}

	gr, err := newGrcron(defaultState, stateFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "not enough arguments")
		os.Exit(1)
	}

	canrun, err := gr.canRun()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if dryRun {
		fmt.Printf("dry-run gr.CurrentState:%s, gr.IsActive:%v finished.\n", gr.CurrentState, canrun)
		return
	}

	if !canrun {
		return
	}

	cmd := exec.Command(args[0])
	for _, arg := range args[1:] {
		cmd.Args = append(cmd.Args, arg)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
