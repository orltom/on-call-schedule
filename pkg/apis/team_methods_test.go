package apis

import (
	"testing"
	"time"
)

func TestVacationDay_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		arg     []byte
		want    *VacationDay
		wantErr bool
	}{
		{
			name:    "Invalid date should throw an error",
			arg:     []byte("abc"),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Parse a valid date",
			arg:     []byte("2024-01-01"),
			want:    &VacationDay{date("2024-01-01")},
			wantErr: false,
		},
		{
			name:    "Parse a valid date in quotes",
			arg:     []byte("\"2024-01-01\""),
			want:    &VacationDay{date("2024-01-01")},
			wantErr: false,
		},
		{
			name:    "parse empty json string",
			arg:     []byte("\"\""),
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &VacationDay{
				Time: time.Now(),
			}
			err := v.UnmarshalJSON(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !(v != nil || tt.want != nil) && tt.want.Equal(v.Time) {
				t.Errorf("time result = %v, want %v", v, tt.want)
			}
		})
	}
}

func date(day string) time.Time {
	t, _ := time.Parse("2006-01-02", day)
	return t
}
