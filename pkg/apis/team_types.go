package apis

import (
	"time"
)

type EmployeeID string

type Team struct {
	Employees []Employee `json:"employees"`
}

type VacationDay struct {
	time.Time
}

type Employee struct {
	ID           EmployeeID    `json:"id"`
	Name         string        `json:"name"`
	VacationDays []VacationDay `json:"vacationDays,omitempty"`
}
