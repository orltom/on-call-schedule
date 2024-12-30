package apis

import (
	"io"
	"time"
)

type Shift struct {
	Start     time.Time
	End       time.Time
	Primary   EmployeeID
	Secondary EmployeeID
}

type Rule interface {
	Match(primary Employee, shifts []Shift, start time.Time, end time.Time) bool
}

type Exporter interface {
	Write(shifts []Shift, writer io.Writer) error
}
