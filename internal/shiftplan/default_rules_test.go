package shiftplan

import (
	"testing"
	"time"

	"github.com/orltom/on-call-schedule/pkg/apis"
)

func TestVacationConflict(t *testing.T) {
	type args struct {
		employee apis.Employee
		start    time.Time
		end      time.Time
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Should detect conflict if employee has holiday between scheduled shift",
			args: args{
				employee: apis.Employee{ID: "a", Name: "a", VacationDays: []string{"2024-01-01"}},
				start:    date("2024-01-01"),
				end:      date("2024-01-02"),
			},
			want: true,
		},
		{
			name: "Should not detect conflict if employee has no holidays",
			args: args{
				employee: apis.Employee{ID: "a", Name: "a", VacationDays: nil},
				start:    date("2024-01-01"),
				end:      date("2024-01-02"),
			},
			want: false,
		},
		{
			name: "Should not detect conflict if employee has holiday outside scheduled shift",
			args: args{
				employee: apis.Employee{ID: "a", Name: "a", VacationDays: []string{"2024-01-03"}},
				start:    date("2024-01-01"),
				end:      date("2024-01-02"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := VacationConflict()
			if got := d.Match(tt.args.employee, nil, tt.args.start, tt.args.end); got != tt.want {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInvolvedInLastSift(t *testing.T) {
	type args struct {
		employee apis.Employee
		shifts   []apis.Shift
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Should not detect conflict when start with the first schedule shift",
			args: args{
				employee: apis.Employee{ID: "a", Name: "a", VacationDays: nil},
				shifts:   nil,
			},
			want: false,
		},
		{
			name: "Should not detect conflict if employee was not on last shift",
			args: args{
				employee: apis.Employee{ID: "a", Name: "a", VacationDays: nil},
				shifts: []apis.Shift{{
					Primary:   "b",
					Secondary: "c",
				}},
			},
			want: false,
		},
		{
			name: "Should detect conflict if employee was on last shift as primary",
			args: args{
				employee: apis.Employee{ID: "a", Name: "a", VacationDays: nil},
				shifts: []apis.Shift{{
					Primary:   "a",
					Secondary: "b",
				}},
			},
			want: true,
		},
		{
			name: "Should detect conflict if employee was on last shift as secondary",
			args: args{
				employee: apis.Employee{ID: "a", Name: "a", VacationDays: nil},
				shifts: []apis.Shift{{
					Primary:   "c",
					Secondary: "a",
				}},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := InvolvedInLastSift()
			if got := d.Match(tt.args.employee, tt.args.shifts, time.Now(), time.Now()); got != tt.want {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}
