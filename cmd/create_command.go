package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/orltom/on-call-schedule/internal/cli"
	"github.com/orltom/on-call-schedule/internal/export"
	"github.com/orltom/on-call-schedule/internal/shiftplan"
	"github.com/orltom/on-call-schedule/pkg/apis"
)

type Format int

const (
	CVS Format = iota
	Table
	JSON
)

func RunCreateShiftPlan(arguments []string) error {
	enums := map[string]Format{"CVS": CVS, "Table": Table, "json": JSON}
	exporters := map[Format]func() apis.Exporter{
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

	start := new(time.Time)
	end := new(time.Time)
	primaryRules := new(string)
	secondaryRules := new(string)
	duration := new(int)
	teamFilePath := new(string)
	outputFormat := Table
	var showHelp bool

	createCommand := flag.NewFlagSet("create", flag.ExitOnError)
	createCommand.BoolVar(&showHelp, "h", false, "help for ocsctl create")
	createCommand.IntVar(duration, "duration", 7*24, "shift duration in hours")
	createCommand.Func("start", "(required) start time of the schedule plan", cli.TimeValueVar(start))
	createCommand.Func("end", "(required) end time of the schedule plan", cli.TimeValueVar(end))
	createCommand.Func("team-file", "(required) path to the file that contain all on-call duties", cli.FilePathVar(teamFilePath))
	createCommand.Func("output", "output format. One of (cvs, table, json)", cli.EnumValueVar(enums, &outputFormat))
	createCommand.StringVar(primaryRules, "primary-rules", "vacation", "Rule to decide which employee should be on-call for the next shift")
	createCommand.StringVar(secondaryRules, "secondary-rules", "vacation", "Rule to decide which employee should be on-call for the next shift")
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

	// validate CLI arguments...
	if ok, missed := cli.RequiredFlagPassed(createCommand, "start", "end", "team-file"); !ok {
		createCommand.Usage()
		return fmt.Errorf("missing required flag: %s", strings.Join(missed, ","))
	}

	// initialize...
	team, err := readTeamFile(*teamFilePath)
	if err != nil {
		return fmt.Errorf("could not read %s: %w", *teamFilePath, err)
	}

	pRules, err := mapPrimaryRules(*primaryRules)
	if err != nil {
		return fmt.Errorf("invalid rules: %w", err)
	}

	sRules, err := mapSecondaryRules(*secondaryRules)
	if err != nil {
		return fmt.Errorf("invalid rules: %w", err)
	}

	// run...
	plan, err := shiftplan.NewShiftPlanner(team.Employees, pRules, sRules).Plan(*start, *end, time.Duration(*duration)*time.Hour)
	if err != nil {
		return fmt.Errorf("can not create on-call schedule: %w", err)
	}

	if err := exporters[outputFormat]().Write(plan, os.Stdout); err != nil {
		return fmt.Errorf("unexpecting error: %w", err)
	}

	return nil
}

func readTeamFile(path string) (apis.Team, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return apis.Team{}, fmt.Errorf("could not read file '%s': %w", path, err)
	}

	var team apis.Team
	if err := json.Unmarshal(content, &team); err != nil {
		return apis.Team{}, fmt.Errorf("could not parse json file '%s': %w", path, err)
	}

	return team, nil
}

func mapPrimaryRules(value string) ([]apis.Rule, error) {
	var rules []apis.Rule
	for _, s := range strings.Split(value, ",") {
		switch strings.ToLower(s) {
		case "vacation":
			rules = append(rules, shiftplan.NewNoVacationOverlap())
		case "minimumoneshiftgap":
			rules = append(rules, shiftplan.NewMinimumPrimaryGapBetweenShifts(1))
		case "minimumtwoshiftgap":
			rules = append(rules, shiftplan.NewMinimumPrimaryGapBetweenShifts(2))
		case "minimumthreeshiftgap":
			rules = append(rules, shiftplan.NewMinimumPrimaryGapBetweenShifts(3))
		case "minimumfourshiftgap":
			rules = append(rules, shiftplan.NewMinimumPrimaryGapBetweenShifts(4))
		default:
			return nil, fmt.Errorf("unknow rule: %s", s)
		}
	}
	return rules, nil
}

func mapSecondaryRules(value string) ([]apis.Rule, error) {
	var rules []apis.Rule
	for _, s := range strings.Split(value, ",") {
		switch strings.ToLower(s) {
		case "vacation":
			rules = append(rules, shiftplan.NewNoVacationOverlap())
		case "minimumoneshiftgap":
			rules = append(rules, shiftplan.NewMinimumSecondaryGapBetweenShifts(1))
		case "minimumtwoshiftgap":
			rules = append(rules, shiftplan.NewMinimumSecondaryGapBetweenShifts(2))
		case "minimumthreeshiftgap":
			rules = append(rules, shiftplan.NewMinimumSecondaryGapBetweenShifts(3))
		case "minimumfourshiftgap":
			rules = append(rules, shiftplan.NewMinimumSecondaryGapBetweenShifts(4))
		default:
			return nil, fmt.Errorf("unknow rule: %s", s)
		}
	}
	return rules, nil
}
