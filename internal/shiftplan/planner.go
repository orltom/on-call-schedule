package shiftplan

import (
	"slices"
	"time"

	"github.com/orltom/on-call-schedule/pkg/apis"
)

type ShiftPlanner struct {
	primaryConflictCheckers  []apis.Rule
	secondaryConflictChecker []apis.Rule
	team                     []apis.Employee
}

func NewDefaultShiftPlanner(team []apis.Employee) *ShiftPlanner {
	return NewShiftPlanner(
		team,
		[]apis.Rule{VacationConflict(), InvolvedInLastSift()},
		[]apis.Rule{VacationConflict(), InvolvedInLastSift()},
	)
}

func NewShiftPlanner(team []apis.Employee, primaryConflictCheckers []apis.Rule, secondaryConflictCheckers []apis.Rule) *ShiftPlanner {
	return &ShiftPlanner{
		team:                     team,
		primaryConflictCheckers:  primaryConflictCheckers,
		secondaryConflictChecker: secondaryConflictCheckers,
	}
}

func (p *ShiftPlanner) Plan(start time.Time, end time.Time, rotation time.Duration) []apis.Shift {
	var plan []apis.Shift

	for s := start; s.Before(end); s = s.Add(rotation) {
		e := s.Add(rotation)

		primary := p.findPrimary(plan, s, e)
		p.team = remove(p.team, primary)

		secondary := p.findSecondary(plan, s, e)
		p.team = remove(p.team, secondary)

		shift := apis.Shift{Start: s, End: e, Primary: primary.ID, Secondary: secondary.ID}
		plan = append(plan, shift)

		p.team = append(p.team, secondary, primary)
	}

	return plan
}

func (p *ShiftPlanner) findPrimary(shifts []apis.Shift, start time.Time, end time.Time) apis.Employee {
	return p.find(shifts, start, end, p.primaryConflictCheckers)
}

func (p *ShiftPlanner) findSecondary(shifts []apis.Shift, start time.Time, end time.Time) apis.Employee {
	return p.find(shifts, start, end, p.secondaryConflictChecker)
}

func (p *ShiftPlanner) find(shifts []apis.Shift, start time.Time, end time.Time, checkers []apis.Rule) apis.Employee {
	for idx := range p.team {
		if p.hasConflict(p.team[idx], checkers, shifts, start, end) {
			continue
		}

		return p.team[idx]
	}

	return p.team[0]
}

func (p *ShiftPlanner) hasConflict(e apis.Employee, conflictCheckers []apis.Rule, shifts []apis.Shift, start time.Time, end time.Time) bool {
	for _, checker := range conflictCheckers {
		if checker.Match(e, shifts, start, end) {
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
