package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"slices"
)

var ErrCommandNotFound = errors.New("command not found")

type command struct {
	name        string
	description string
	run         func([]string) error
}

func main() {
	commands := []command{
		{
			name:        "create",
			description: "create on-call shift plan",
			run:         RunCreateShiftPlan,
		},
	}
	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, `ocs helps to create and sync on-call shifts

Usage:
  ocs [command]

Available Commands:
`)
		for _, c := range commands {
			_, _ = fmt.Fprintf(os.Stderr, "  %s\t%s\n", c.name, c.description)
		}
	}
	flag.Parse()

	if len(flag.Args()) < 1 {
		flag.Usage()
		_, _ = fmt.Fprintf(os.Stderr, "\n\nError: missing command")
		os.Exit(128)
	}

	cmd := flag.Args()[0]
	if err := runCommand(cmd, commands); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "\n\nError: %s", err.Error())
		os.Exit(128)
	}
}

func runCommand(cmd string, commands []command) error {
	cmdIdx := slices.IndexFunc(commands, func(c command) bool {
		return c.name == cmd
	})
	if cmdIdx < 0 {
		return fmt.Errorf("%w: %s", ErrCommandNotFound, cmd)
	}

	if err := commands[cmdIdx].run(os.Args[2:]); err != nil {
		return fmt.Errorf("\n\n%w", err)
	}

	return nil
}
