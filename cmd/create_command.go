package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/orltom/on-call-schedule/internal/cli"
	"github.com/orltom/on-call-schedule/internal/export"
	"github.com/orltom/on-call-schedule/internal/shiftplan"
	"github.com/orltom/on-call-schedule/pkg/apis"
)

var (
	ErrMissingArguments = errors.New("missing required flag")
	ErrInvalidArgument  = errors.New("invalid value")
)

type Converter int

const (
	CVS Converter = iota
	Table
)

func RunCreateShiftPlan(arguments []string) error {
	enums := map[string]Converter{"CVS": CVS, "Table": Table}
	converters := map[Converter]func() apis.Exporter{
		Table: func() apis.Exporter {
			return export.NewTableExporter()
		},
		CVS: func() apis.Exporter {
			return export.NewCVSCExporter()
		},
	}

	str := new(time.Time)
	end := new(time.Time)
	duration := new(int)
	teamFilePath := new(string)
	transform := Table

	createCommand := flag.NewFlagSet("create", flag.ExitOnError)
	createCommand.IntVar(duration, "rotation", 7*24, "rotation time in hours.")
	createCommand.Func("start", "(required) start time of the schedule plan", cli.TimeValueVar(str))
	createCommand.Func("end", "(required) end time of the schedule plan", cli.TimeValueVar(end))
	createCommand.Func("team-file", "(required) path to the file that contain all on-call duties", cli.FilePathVar(teamFilePath))
	createCommand.Func("output", "output format. One of (cvs, table)", cli.EnumValueVar(enums, &transform))
	createCommand.Usage = func() {
		fmt.Fprintf(os.Stderr, `Create on-call shift plan

Usage:
  ocs create [flags]

Flags:
`)
		createCommand.PrintDefaults()
	}

	if err := createCommand.Parse(arguments); err != nil {
		createCommand.Usage()
		return fmt.Errorf("could not parse CLI arguments: %w", err)
	}

	// check that the required flags are set...
	if !cli.IsFlagPassed(createCommand, "start") {
		createCommand.Usage()
		return fmt.Errorf("%w: %s", ErrMissingArguments, "start")
	}

	if !cli.IsFlagPassed(createCommand, "end") {
		createCommand.Usage()
		return fmt.Errorf("%w: %s", ErrMissingArguments, "end")
	}

	if !cli.IsFlagPassed(createCommand, "team-file") {
		createCommand.Usage()
		return fmt.Errorf("%w: %s", ErrMissingArguments, "team-file")
	}

	// validate input data...
	if str == end {
		createCommand.Usage()
		return ErrInvalidArgument
	}

	// initialize and run...
	team, err := parse(*teamFilePath)
	if err != nil {
		return fmt.Errorf("%w: invalid team config file: %w", ErrInvalidArgument, err)
	}

	plan := shiftplan.NewDefaultShiftPlanner(team.Employees).Plan(*str, *end, time.Duration(*duration)*time.Hour)
	if err := converters[transform]().Write(plan, os.Stdout); err != nil {
		return fmt.Errorf("unexpecting error: %w", err)
	}

	return nil
}

func parse(path string) (apis.Team, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return apis.Team{}, fmt.Errorf("could not read file '%s': %w", path, err)
	}

	var team apis.Team
	if err := json.Unmarshal(content, &team); err != nil {
		return apis.Team{}, fmt.Errorf("could not pars json file '%s': %w", path, err)
	}

	return team, nil
}
