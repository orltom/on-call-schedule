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

func NewMinimumPrimaryGapBetweenShifts(gap int) apis.RuleFunc {
	return func(e apis.Employee, shifts []apis.Shift, _ time.Time, _ time.Time) bool {
		if len(shifts) == 0 {
			return false
		}

		lowIndex := 0
		if len(shifts) > gap {
			lowIndex = len(shifts) - gap
		}
		for idx := range shifts[lowIndex:] {
			if shifts[idx].Primary == e.ID {
				return true
			}
		}
		return false
	}
}

func NewMinimumSecondaryGapBetweenShifts(gap int) apis.RuleFunc {
	return func(e apis.Employee, shifts []apis.Shift, _ time.Time, _ time.Time) bool {
		if len(shifts) == 0 {
			return false
		}

		lowIndex := 0
		if len(shifts) > gap {
			lowIndex = len(shifts) - gap
		}
		for idx := range shifts[lowIndex:] {
			if shifts[idx].Secondary == e.ID {
				return true
			}
		}
		return false
	}
}
