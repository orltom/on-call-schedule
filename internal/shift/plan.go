package shift

import (
	"slices"
	"time"

	"github.com/orltom/on-call-schedule/pkg/apis"
)

type ShiftRotation struct {
	primaryConflictCheckers  []apis.ShiftConflictChecker
	secondaryConflictChecker []apis.ShiftConflictChecker
	team                     []apis.Employee
}

func NewDefaultShiftRotation(team []apis.Employee) *ShiftRotation {
	return &ShiftRotation{
		team:                     team,
		primaryConflictCheckers:  []apis.ShiftConflictChecker{VacationConflict, InvolvedInLastSift},
		secondaryConflictChecker: []apis.ShiftConflictChecker{VacationConflict, InvolvedInLastSift},
	}
}

func NewShiftRotation(team []apis.Employee, primaryConflictCheckers []apis.ShiftConflictChecker, secondaryConflictCheckers []apis.ShiftConflictChecker) *ShiftRotation {
	return &ShiftRotation{
		team:                     team,
		primaryConflictCheckers:  primaryConflictCheckers,
		secondaryConflictChecker: secondaryConflictCheckers,
	}
}

func (sr *ShiftRotation) Plan(start time.Time, end time.Time, rotation time.Duration) []apis.Shift {
	var shiftplan []apis.Shift

	for s := start; s.Before(end); s = s.Add(rotation) {
		e := s.Add(rotation)

		primary := sr.findPrimary(shiftplan, s, e)
		sr.team = remove(sr.team, primary)

		secondary := sr.findSecondary(shiftplan, s, e)
		sr.team = remove(sr.team, secondary)

		shift := apis.Shift{Start: s, End: e, Primary: primary, Secondary: secondary}
		shiftplan = append(shiftplan, shift)

		sr.team = append(sr.team, secondary, primary)
	}

	return shiftplan
}

func (sr *ShiftRotation) findPrimary(shifts []apis.Shift, start time.Time, end time.Time) apis.Employee {
	return sr.find(shifts, start, end, sr.primaryConflictCheckers)
}

func (sr *ShiftRotation) findSecondary(shifts []apis.Shift, start time.Time, end time.Time) apis.Employee {
	return sr.find(shifts, start, end, sr.secondaryConflictChecker)
}

func (sr *ShiftRotation) find(shifts []apis.Shift, start time.Time, end time.Time, checkers []apis.ShiftConflictChecker) apis.Employee {
	for idx := range sr.team {
		if sr.hasConflict(sr.team[idx], checkers, shifts, start, end) {
			continue
		}

		return sr.team[idx]
	}

	return sr.team[0]
}

func (sr *ShiftRotation) hasConflict(e apis.Employee, conflictCheckers []apis.ShiftConflictChecker, shifts []apis.Shift, start time.Time, end time.Time) bool {
	for _, checker := range conflictCheckers {
		if checker(e, shifts, start, end) {
			return true
		}
	}

	return false
}

func remove(list []apis.Employee, e apis.Employee) []apis.Employee {
	return slices.DeleteFunc(list, func(employee apis.Employee) bool {
		return e.Name == employee.Name
	})
}
