package apis

import (
	"io"
	"time"
)

type Shift struct {
	Start     time.Time  `json:"start"`
	End       time.Time  `json:"end"`
	Primary   EmployeeID `json:"primary"`
	Secondary EmployeeID `json:"secondary"`
}

type Rule interface {
	Match(employee Employee, shifts []Shift, start time.Time, end time.Time) bool
}

type RuleFunc func(Employee, []Shift, time.Time, time.Time) bool

func (r RuleFunc) Match(employee Employee, shifts []Shift, start time.Time, end time.Time) bool {
	return r(employee, shifts, start, end)
}

type Exporter interface {
	Write(plan []Shift, writer io.Writer) error
}
