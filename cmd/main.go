package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"slices"
)

var ErrCommandNotFound = errors.New("command not found")

type command struct {
	name        string
	description string
	run         func(io.Writer, []string) error
}

func main() {
	commands := []command{
		{
			name:        "create",
			description: "create on-call shift plan",
			run:         RunCreateShiftPlan,
		},
	}

	var showHelp bool
	flag.BoolVar(&showHelp, "h", false, "help for ocsctl")
	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "%s helps to create and sync on-call shifts\n", os.Args[0])
		fmt.Fprintf(os.Stdout, "Usage:\n")
		fmt.Fprintf(os.Stdout, "  %s\t[command] [flags]\n", os.Args[0])
		fmt.Fprintf(os.Stdout, "\nAvailable Commands:\n")
		for idx := range commands {
			fmt.Fprintf(os.Stdout, "  %s    %s\n", commands[idx].name, commands[idx].description)
		}
		fmt.Fprintf(os.Stdout, "\nFlags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stdout, "\nUse \"%s [command] -h\" for more information about a command\n", os.Args[0])
	}
	flag.Parse()

	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if len(flag.Args()) < 1 {
		flag.Usage()
		_, _ = fmt.Fprintf(os.Stderr, "\n\nError: missing command\n")
		os.Exit(128)
	}

	cmd := flag.Args()[0]
	if err := runCommand(cmd, commands); err != nil {
		fmt.Fprintf(os.Stderr, "\n\nError: %s\n", err.Error())
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

	if err := commands[cmdIdx].run(os.Stdout, os.Args[2:]); err != nil {
		return fmt.Errorf("\n\n%w", err)
	}

	return nil
}
