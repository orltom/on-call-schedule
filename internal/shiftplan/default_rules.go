package shiftplan

import (
	"time"

	"github.com/orltom/on-call-schedule/pkg/apis"
)

var _ apis.Rule = VacationConflict()

var _ apis.Rule = InvolvedInLastSift()

type DefaultRule struct {
	fn func(e apis.Employee, _ []apis.Shift, start time.Time, end time.Time) bool
}

func (d *DefaultRule) Match(employee apis.Employee, shifts []apis.Shift, start time.Time, end time.Time) bool {
	return d.fn(employee, shifts, start, end)
}

func VacationConflict() *DefaultRule {
	return &DefaultRule{
		fn: func(e apis.Employee, _ []apis.Shift, start time.Time, end time.Time) bool {
			days := e.VacationDays
			for vIdx := range days {
				day := days[vIdx]
				if (day.After(start) && day.Before(end)) || day.Equal(start) || day.Equal(end) {
					return true
				}
			}
			return false
		},
	}
}

func InvolvedInLastSift() *DefaultRule {
	return &DefaultRule{
		fn: func(e apis.Employee, shifts []apis.Shift, _ time.Time, _ time.Time) bool {
			if len(shifts) == 0 {
				return false
			}
			lastShift := shifts[len(shifts)-1]
			if lastShift.Primary == e.ID || lastShift.Secondary == e.ID {
				return true
			}
			return false
		},
	}
}
