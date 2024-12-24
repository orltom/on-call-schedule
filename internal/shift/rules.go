package shift

import (
	"time"

	"github.com/orltom/on-call-schedule/pkg/apis"
)

var _ apis.ShiftConflictChecker = VacationConflict

var _ apis.ShiftConflictChecker = InvolvedInLastSift

func VacationConflict(e apis.Employee, _ []apis.Shift, start time.Time, end time.Time) bool {
	for vIdx := range e.VacationDays {
		day := e.VacationDays[vIdx]
		if (day.After(start) && day.Before(end)) || day.Equal(start) || day.Equal(end) {
			return true
		}
	}

	return false
}

func InvolvedInLastSift(e apis.Employee, shifts []apis.Shift, _ time.Time, _ time.Time) bool {
	if len(shifts) == 0 {
		return false
	}

	lastShift := shifts[len(shifts)-1]
	if lastShift.Primary.Name == e.Name || lastShift.Secondary.Name == e.Name {
		return true
	}

	return false
}
