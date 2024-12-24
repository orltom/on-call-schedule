package apis

import "time"

type Config struct {
	Employees []Employee `json:"employees"`
}

type Employee struct {
	Name         string      `json:"name"`
	VacationDays []time.Time `json:"vacationDays,omitempty"`
}
