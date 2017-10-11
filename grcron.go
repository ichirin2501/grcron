package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"syscall"
)

const Version string = "0.0.1"

type Grcron struct {
	StateFile    string
	DefaultState string
	CurrentState string
}

func (gr Grcron) Validate() error {
	_, err := os.Stat(gr.StateFile)
	if err != nil {
		return err
	}
	if !(gr.DefaultState == "active" || gr.DefaultState == "passive") {
		return fmt.Errorf("The Value of DefaultState:%s is incorrect.", gr.DefaultState)
	}
	return nil
}

func (gr *Grcron) ParseState() error {
	if gr == nil {
		return fmt.Errorf("Don't run nil Pointer Receiver.")
	}
	f, err := os.Open(gr.StateFile)
	defer f.Close()
	if err != nil {
		return err
	}

	sc := bufio.NewScanner(f)
	if !sc.Scan() {
		return sc.Err()
	}
	st := sc.Text()
	switch st {
	case "active", "passive":
		gr.CurrentState = st
	default:
		fmt.Fprintf(os.Stderr, "corrupted state file('%s') (content='%s'), staying at gr.DefaultState('%s')\n", gr.StateFile, st, gr.DefaultState)
		gr.CurrentState = gr.DefaultState
	}
	return nil
}
func (gr Grcron) IsActive() (bool, error) {
	cmd := exec.Command("sh", "-c", "ps cax | grep -q keepalived")
	err := cmd.Run()

	// 異常終了はkeepalivedプロセスがいないとみなす
	if _, ok := err.(*exec.ExitError); ok {
		return false, fmt.Errorf("keepalived is probably down.")
	}

	return gr.CurrentState == "active", nil
}

func main() {
	var (
		showVersion bool
		dryRun      bool
	)

	gr := &Grcron{}
	flag.StringVar(&gr.StateFile, "f", "/var/run/grcron/state", "grcron state file.")
	flag.StringVar(&gr.DefaultState, "s", "passive", "grcron default state.")
	flag.BoolVar(&showVersion, "version", false, "show version number.")
	flag.BoolVar(&showVersion, "v", false, "show version number.")
	flag.BoolVar(&dryRun, "dryrun", false, "dry-run.")
	flag.BoolVar(&dryRun, "n", false, "dry-run.")
	flag.Parse()
	args := flag.Args()

	if showVersion {
		fmt.Printf("grcron %s, %s built for %s/%s\n", Version, runtime.Version(), runtime.GOOS, runtime.GOARCH)
		return
	}

	if err := gr.Validate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := gr.ParseState(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "not enough arguments")
		os.Exit(1)
	}

	isa, err := gr.IsActive()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if dryRun {
		fmt.Printf("dry-run gr.CurrentState:%s, gr.IsActive:%v finished.\n", gr.CurrentState, isa)
		return
	}

	if !isa {
		return
	}

	// run !!
	binary, err := exec.LookPath(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := syscall.Exec(binary, args, os.Environ()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
