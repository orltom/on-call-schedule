package shiftplan

import (
	"fmt"
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
		[]apis.Rule{NewNoVacationOverlap()},
		[]apis.Rule{NewNoVacationOverlap()},
	)
}

func NewShiftPlanner(team []apis.Employee, primaryConflictCheckers []apis.Rule, secondaryConflictCheckers []apis.Rule) *ShiftPlanner {
	t := make([]apis.Employee, len(team))
	copy(t, team)

	p := make([]apis.Rule, len(primaryConflictCheckers))
	copy(p, primaryConflictCheckers)

	s := make([]apis.Rule, len(secondaryConflictCheckers))
	copy(s, secondaryConflictCheckers)

	return &ShiftPlanner{
		team:                     t,
		primaryConflictCheckers:  p,
		secondaryConflictChecker: s,
	}
}

func (p *ShiftPlanner) Plan(start time.Time, end time.Time, duration time.Duration) ([]apis.Shift, error) {
	var plan []apis.Shift

	for s := start; s.Before(end); s = s.Add(duration) {
		e := s.Add(duration)

		primary, err := p.findPrimary(plan, s, e)
		if err != nil {
			return nil, err
		}
		p.team = remove(p.team, primary)

		secondary, err := p.findSecondary(plan, s, e)
		if err != nil {
			return nil, err
		}
		p.team = remove(p.team, secondary)

		shift := apis.Shift{Start: s, End: e, Primary: primary.ID, Secondary: secondary.ID}
		plan = append(plan, shift)

		p.team = append(p.team, secondary, primary)
	}

	return plan, nil
}

func (p *ShiftPlanner) findPrimary(shifts []apis.Shift, start time.Time, end time.Time) (apis.Employee, error) {
	return p.find(shifts, start, end, p.primaryConflictCheckers)
}

func (p *ShiftPlanner) findSecondary(shifts []apis.Shift, start time.Time, end time.Time) (apis.Employee, error) {
	return p.find(shifts, start, end, p.secondaryConflictChecker)
}

func (p *ShiftPlanner) find(shifts []apis.Shift, start time.Time, end time.Time, checkers []apis.Rule) (apis.Employee, error) {
	for idx := range p.team {
		if p.hasConflict(p.team[idx], checkers, shifts, start, end) {
			continue
		}

		return p.team[idx], nil
	}

	return apis.Employee{}, fmt.Errorf("could not find available duty between %s and %s", start, end)
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
