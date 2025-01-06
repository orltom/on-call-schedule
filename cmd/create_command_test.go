package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunCreateShiftPlan(t *testing.T) {
	dir := t.TempDir()
	tempFilePath := filepath.Join(dir, "team.json")
	content := `{"employees": [{"id": "joe@example.com", "name": "Joe"}, {"id": "jan@example.com", "name": "Jan"}]}`
	if err := os.WriteFile(tempFilePath, []byte(content), 0o600); err != nil {
		t.FailNow()
	}

	type args struct {
		arguments []string
	}
	tests := []struct {
		name       string
		args       args
		wantWriter string
		wantErrMsg string
	}{
		{
			name:       "Print Usage when required flags are missing",
			wantWriter: "Usage",
			wantErrMsg: "missing required flag: start,end,team-file",
		},
		{
			name: "Print Usage when use help flag",
			args: args{
				arguments: []string{
					"-h",
				},
			},
			wantWriter: "Usage",
			wantErrMsg: "",
		},
		{
			name: "Print an on-call schedule in JSON format",
			args: args{
				arguments: []string{
					"--start", "2024-01-01 00:00:00",
					"--end", "2024-01-08 00:00:00",
					"--team-file", tempFilePath,
					"--output", "json",
				},
			},
			wantWriter: `[{"start":"2024-01-01T00:00:00Z","end":"2024-01-08T00:00:00Z","primary":"joe@example.com","secondary":"jan@example.com"}]`,
			wantErrMsg: "",
		},
		{
			name: "Print an error if shift plan cannot be created if rules are too restricted",
			args: args{
				arguments: []string{
					"--start", "2024-01-01 00:00:00",
					"--end", "2024-01-15 00:00:00",
					"--team-file", tempFilePath,
					"--primary-rules", "minimumfourshiftgap",
					"--output", "json",
				},
			},
			wantWriter: "",
			wantErrMsg: "can not create on-call schedule: could not find available duty between 2024-01-08 00:00:00 +0000 UTC and 2024-01-15 00:00:00 +0000 UTC",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			err := RunCreateShiftPlan(writer, tt.args.arguments)
			if (err != nil) && tt.wantErrMsg != err.Error() {
				t.Errorf("RunCreateShiftPlan() error = %v, wantErr %v", err, tt.wantErrMsg)
				return
			}
			if gotWriter := writer.String(); !strings.Contains(gotWriter, tt.wantWriter) {
				t.Errorf("RunCreateShiftPlan() gotWriter = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}
