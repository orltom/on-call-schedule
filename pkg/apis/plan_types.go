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
	Match(primary Employee, shifts []Shift, start time.Time, end time.Time) bool
}

type Exporter interface {
	Write(plan []Shift, writer io.Writer) error
}
