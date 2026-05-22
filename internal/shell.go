package internal

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	ShellEnvVar        = "SHELL"
	DefaultShell       = "/bin/bash"
	SubShellVar        = "WHOIAM_SUBSHELL"
	SubShellVarEnabled = "WHOIAM_SUBSHELL=1"
)

type SubShell struct {
	CommandArgs []string
}

func GetShell() string {
	shell := os.Getenv(ShellEnvVar)
	if shell == "" {
		shell = DefaultShell
	}
	return shell
}

func GetIsWhoIAmSubShell() bool {
	return os.Getenv(SubShellVar) == "1"
}

func NewSubShell(args ...string) (*exec.Cmd, error) {
	var cmd *exec.Cmd

	if GetIsWhoIAmSubShell() {
		return nil, fmt.Errorf("cannot create a subshell from a subshell. Exit the current subshell by typing exit and try again")
	}

	shell := GetShell()
	if len(args) > 0 {
		cmd = exec.Command(shell, "-c", strings.Join(args, " "))
	} else {
		cmd = exec.Command(shell)
	}

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	cmd.Env = append(os.Environ(), SubShellVarEnabled)
	return cmd, nil
}
