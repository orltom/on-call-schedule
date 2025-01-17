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
		duration time.Duration
	}
	tests := []struct {
		name      string
		employees []apis.Employee
		args      args
		want      []apis.Shift
		wantErr   bool
	}{
		{
			name: "Without declared holidays, a daily schedule should be generated according to the order of the employees",
			employees: []apis.Employee{
				{ID: "a@test.ch", Name: "a", VacationDays: nil},
				{ID: "b@test.ch", Name: "b", VacationDays: nil},
			},
			args: args{
				start:    date("2020-04-01"),
				end:      date("2020-04-04"),
				duration: 24 * time.Hour,
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
			wantErr: false,
		},
		{
			name: "Should not schedule Primary on holiday days",
			employees: []apis.Employee{
				{ID: "a@test.ch", Name: "a", VacationDays: vacationDays("2020-04-01")},
				{ID: "b@test.ch", Name: "b", VacationDays: nil},
				{ID: "c@test.ch", Name: "c", VacationDays: nil},
				{ID: "d@test.ch", Name: "d", VacationDays: nil},
			},
			args: args{
				start:    date("2020-04-01"),
				end:      date("2020-04-05"),
				duration: 24 * time.Hour,
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
			wantErr: false,
		},
		{
			name: "Should not schedule Secondary on holiday days",
			employees: []apis.Employee{
				{ID: "a@test.ch", Name: "a", VacationDays: nil},
				{ID: "b@test.ch", Name: "b", VacationDays: vacationDays("2020-04-01")},
				{ID: "c@test.ch", Name: "c", VacationDays: nil},
				{ID: "d@test.ch", Name: "d", VacationDays: nil},
			},
			args: args{
				start:    date("2020-04-01"),
				end:      date("2020-04-05"),
				duration: 24 * time.Hour,
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
			wantErr: false,
		},
		{
			name: "Should return an error if can not find next available duty",
			employees: []apis.Employee{
				{ID: "a@test.ch", Name: "a", VacationDays: vacationDays("2020-04-01")},
				{ID: "b@test.ch", Name: "b", VacationDays: vacationDays("2020-04-01")},
				{ID: "c@test.ch", Name: "c", VacationDays: vacationDays("2020-04-01")},
				{ID: "d@test.ch", Name: "d", VacationDays: vacationDays("2020-04-01")},
			},
			args: args{
				start:    date("2020-04-01"),
				end:      date("2020-04-04"),
				duration: 24 * time.Hour,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should return an error if can not find next secondary",
			employees: []apis.Employee{
				{ID: "a@test.ch", Name: "a", VacationDays: nil},
				{ID: "b@test.ch", Name: "b", VacationDays: vacationDays("2020-04-01")},
				{ID: "c@test.ch", Name: "c", VacationDays: vacationDays("2020-04-01")},
				{ID: "d@test.ch", Name: "d", VacationDays: vacationDays("2020-04-01")},
			},
			args: args{
				start:    date("2020-04-01"),
				end:      date("2020-04-04"),
				duration: 24 * time.Hour,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			planner := NewDefaultShiftPlanner(tt.employees)
			got, err := planner.Plan(tt.args.start, tt.args.end, tt.args.duration)
			if (err != nil) != tt.wantErr {
				t.Errorf("Plan() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
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

func vacationDays(days ...string) []apis.VacationDay {
	var tmp = make([]apis.VacationDay, len(days))
	for _, day := range days {
		tmp = append(tmp, apis.VacationDay{Time: date(day)})
	}

	return tmp
}
