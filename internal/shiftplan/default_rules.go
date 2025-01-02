package shiftplan

import (
	"time"

	"github.com/orltom/on-call-schedule/pkg/apis"
)

func NewNoVacationOverlap() apis.RuleFunc {
	return func(e apis.Employee, _ []apis.Shift, start time.Time, end time.Time) bool {
		days := e.VacationDays
		for vIdx := range days {
			day := days[vIdx]
			if (day.After(start) && day.Before(end)) || day.Equal(start) || day.Equal(end) {
				return true
			}
		}
		return false
	}
}

func NewMinimumGapBetweenShifts(gap int) apis.RuleFunc {
	return func(e apis.Employee, shifts []apis.Shift, _ time.Time, _ time.Time) bool {
		if len(shifts) == 0 {
			return false
		}
		for idx := range shifts[len(shifts)-gap:] {
			if shifts[idx].Primary == e.ID || shifts[idx].Secondary == e.ID {
				return true
			}
		}
		return false
	}
}
