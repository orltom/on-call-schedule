package apis

import (
	"fmt"
	"time"
)

type Shift struct {
	Start     time.Time
	End       time.Time
	Primary   Employee
	Secondary Employee
}

func (s Shift) String() string {
	return fmt.Sprintf("%s - %s = %s (%s)\n", s.Start.Format("2006-01-02"), s.End.Format("2006-01-02"), s.Primary.Name, s.Secondary.Name)
}

type ShiftConflictChecker func(primary Employee, shifts []Shift, start time.Time, end time.Time) bool
