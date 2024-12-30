package shiftplan

import (
	"reflect"
	"testing"
	"time"

	"github.com/orltom/on-call-schedule/pkg/apis"
)

func TestPlanner_Plan(t *testing.T) {
	type args struct {
		start    time.Time
		end      time.Time
		rotation time.Duration
	}
	tests := []struct {
		name      string
		employees []apis.Employee
		args      args
		want      []apis.Shift
	}{
		{
			name: "Without declared holidays, a daily schedule should be generated according to the order of the employees",
			employees: []apis.Employee{
				{
					ID:           "a@test.ch",
					Name:         "a",
					VacationDays: nil,
				},
				{
					ID:           "b@test.ch",
					Name:         "b",
					VacationDays: nil,
				},
			},
			args: args{
				start:    date("2020-04-01"),
				end:      date("2020-04-04"),
				rotation: 24 * time.Hour,
			},
			want: []apis.Shift{
				{
					Start:     date("2020-04-01"),
					End:       date("2020-04-02"),
					Primary:   "a@test.ch",
					Secondary: "b@test.ch",
				},
				{
					Start:     date("2020-04-02"),
					End:       date("2020-04-03"),
					Primary:   "b@test.ch",
					Secondary: "a@test.ch",
				},
				{
					Start:     date("2020-04-03"),
					End:       date("2020-04-04"),
					Primary:   "a@test.ch",
					Secondary: "b@test.ch",
				},
			},
		},
		{
			name: "Should not schedule Primary on holiday days",
			employees: []apis.Employee{
				{
					ID:           "a@test.ch",
					Name:         "a",
					VacationDays: []string{"2020-04-01"},
				},
				{
					ID:           "b@test.ch",
					Name:         "b",
					VacationDays: nil,
				},
				{
					ID:           "c@test.ch",
					Name:         "c",
					VacationDays: nil,
				},
				{
					ID:           "d@test.ch",
					Name:         "d",
					VacationDays: nil,
				},
			},
			args: args{
				start:    date("2020-04-01"),
				end:      date("2020-04-05"),
				rotation: 24 * time.Hour,
			},
			want: []apis.Shift{
				{
					Start:     date("2020-04-01"),
					End:       date("2020-04-02"),
					Primary:   "b@test.ch",
					Secondary: "c@test.ch",
				},
				{
					Start:     date("2020-04-02"),
					End:       date("2020-04-03"),
					Primary:   "a@test.ch",
					Secondary: "d@test.ch",
				},
				{
					Start:     date("2020-04-03"),
					End:       date("2020-04-04"),
					Primary:   "c@test.ch",
					Secondary: "b@test.ch",
				},
				{
					Start:     date("2020-04-04"),
					End:       date("2020-04-05"),
					Primary:   "d@test.ch",
					Secondary: "a@test.ch",
				},
			},
		},
		{
			name: "Should not schedule Secondary on holiday days",
			employees: []apis.Employee{
				{
					ID:           "a@test.ch",
					Name:         "a",
					VacationDays: nil,
				},
				{
					ID:           "b@test.ch",
					Name:         "b",
					VacationDays: []string{"2020-04-01"},
				},
				{
					ID:           "c@test.ch",
					Name:         "c",
					VacationDays: nil,
				},
				{
					ID:           "d@test.ch",
					Name:         "d",
					VacationDays: nil,
				},
			},
			args: args{
				start:    date("2020-04-01"),
				end:      date("2020-04-05"),
				rotation: 24 * time.Hour,
			},
			want: []apis.Shift{
				{
					Start:     date("2020-04-01"),
					End:       date("2020-04-02"),
					Primary:   "a@test.ch",
					Secondary: "c@test.ch",
				},
				{
					Start:     date("2020-04-02"),
					End:       date("2020-04-03"),
					Primary:   "b@test.ch",
					Secondary: "d@test.ch",
				},
				{
					Start:     date("2020-04-03"),
					End:       date("2020-04-04"),
					Primary:   "c@test.ch",
					Secondary: "a@test.ch",
				},
				{
					Start:     date("2020-04-04"),
					End:       date("2020-04-05"),
					Primary:   "d@test.ch",
					Secondary: "b@test.ch",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			planner := NewDefaultShiftPlanner(tt.employees)
			if got := planner.Plan(tt.args.start, tt.args.end, tt.args.rotation); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Plan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func date(s string) time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		panic(err)
	}

	return t
}
