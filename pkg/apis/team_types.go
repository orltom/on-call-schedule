package apis

import (
	"encoding/json"
	"fmt"
	"time"
)

type EmployeeID string

type Team struct {
	Employees []Employee `json:"employees"`
}

type Employee struct {
	ID              EmployeeID `json:"id"`
	Name            string     `json:"name"`
	VacationDays    []string   `json:"vacationDays,omitempty"`
	parsedVacations []time.Time
}

func (e *Employee) UnmarshalJSON(data []byte) error {
	type Alias Employee
	temp := &struct {
		*Alias
	}{
		Alias: (*Alias)(e),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	parsedDay, err := parseDateOnly(e.VacationDays)
	if err != nil {
		return err
	}
	e.parsedVacations = parsedDay

	return nil
}

func (e *Employee) Vacations() []time.Time {
	if len(e.VacationDays) != len(e.parsedVacations) {
		e.parsedVacations, _ = parseDateOnly(e.VacationDays)
	}
	return e.parsedVacations
}

func parseDateOnly(dateStr []string) ([]time.Time, error) {
	days := make([]time.Time, len(dateStr))
	for idx, day := range dateStr {
		t, err := time.Parse(time.DateOnly, day)
		if err != nil {
			return nil, fmt.Errorf("failed to parse date '%s'", day)
		}
		days[idx] = t
	}
	return days, nil
}
