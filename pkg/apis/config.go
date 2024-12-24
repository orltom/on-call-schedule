package apis

import "time"

type Config struct {
	Employees []Employee `json:"employees"`
}

type EmployeeID string

type Employee struct {
	ID           EmployeeID  `json:"id"`
	Name         string      `json:"name"`
	VacationDays []time.Time `json:"vacationDays,omitempty"`
}
