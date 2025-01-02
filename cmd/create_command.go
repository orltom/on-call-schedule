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
	JSON
)

func RunCreateShiftPlan(arguments []string) error {
	enums := map[string]Converter{"CVS": CVS, "Table": Table, "json": JSON}
	converters := map[Converter]func() apis.Exporter{
		Table: func() apis.Exporter {
			return export.NewTableExporter()
		},
		CVS: func() apis.Exporter {
			return export.NewCVSCExporter()
		},
		JSON: func() apis.Exporter {
			return export.NewJSONExporter()
		},
	}

	str := new(time.Time)
	end := new(time.Time)
	duration := new(int)
	teamFilePath := new(string)
	transform := Table
	var showHelp bool

	createCommand := flag.NewFlagSet("create", flag.ExitOnError)
	createCommand.BoolVar(&showHelp, "h", false, "help for ocsctl create")
	createCommand.IntVar(duration, "duration", 7*24, "shift duration in hours.")
	createCommand.Func("start", "(required) start time of the schedule plan", cli.TimeValueVar(str))
	createCommand.Func("end", "(required) end time of the schedule plan", cli.TimeValueVar(end))
	createCommand.Func("team-file", "(required) path to the file that contain all on-call duties", cli.FilePathVar(teamFilePath))
	createCommand.Func("output", "output format. One of (cvs, table, json)", cli.EnumValueVar(enums, &transform))
	createCommand.Usage = func() {
		fmt.Fprintf(os.Stdout, "Create on-call schedule\n")
		fmt.Fprintf(os.Stdout, "\nUsage\n")
		fmt.Fprintf(os.Stdout, "  %s create [flags]\n", os.Args[0])
		fmt.Fprintf(os.Stdout, "\nFlags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stdout, "\nUse \"%s create -h\" for more information about a command\n", os.Args[0])
	}

	if err := createCommand.Parse(arguments); err != nil {
		createCommand.Usage()
		return fmt.Errorf("could not parse CLI arguments: %w", err)
	}

	if showHelp {
		createCommand.Usage()
		return nil
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

	plan, err := shiftplan.NewDefaultShiftPlanner(team.Employees).Plan(*str, *end, time.Duration(*duration)*time.Hour)
	if err != nil {
		return fmt.Errorf("can not create on-call schedule: %w", err)
	}

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
