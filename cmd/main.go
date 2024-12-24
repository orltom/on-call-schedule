package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/orltom/on-call-schedule/internal/shift"
	"github.com/orltom/on-call-schedule/pkg/apis"
)

func main() {
	rotationFlag := flag.Int("rotation", 24*7, "Rotation time in hours. Default is one week")
	startFlag := flag.String("start", "", "Shift start date")
	endFlag := flag.String("end", "", "Shift end date")
	configFileFlag := flag.String("config", "config.json", "Config file path")
	flag.Parse()

	start, err := time.Parse("2006-01-02", *startFlag)
	if err != nil {
		flag.PrintDefaults()
		fmt.Println("Invalid start time")
		os.Exit(1)
	}

	end, err := time.Parse("2006-01-02", *endFlag)
	if err != nil {
		flag.PrintDefaults()
		fmt.Println("Invalid end time")
		os.Exit(1)
	}

	config, err := os.ReadFile(*configFileFlag)
	if err != nil {
		flag.PrintDefaults()
		fmt.Println("Config file not found or invalid")
		os.Exit(1)
	}

	var c apis.Config
	if err := json.Unmarshal(config, &c); err != nil {
		flag.PrintDefaults()
		fmt.Println("Config file could not be parsed")
		os.Exit(1)
	}

	duration := time.Hour * time.Duration(*rotationFlag)
	shifts := shift.NewDefaultShiftRotation(c.Employees).Plan(start, end, duration)

	fmt.Println(shifts)
}
