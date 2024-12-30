package export

import (
	"strings"
	"testing"

	"github.com/orltom/on-call-schedule/pkg/apis"
)

func TestJSONExporter_Write(t *testing.T) {
	type args struct {
		plan []apis.Shift
	}
	tests := []struct {
		name       string
		args       args
		wantWriter string
		wantErr    bool
	}{
		{
			name: "Should create valid CSV",
			args: args{
				plan: []apis.Shift{
					{
						Start:     date("2024-12-25"),
						End:       date("2024-12-26"),
						Primary:   "a",
						Secondary: "b",
					},
					{
						Start:     date("2024-12-26"),
						End:       date("2024-12-27"),
						Primary:   "c",
						Secondary: "d",
					},
				},
			},
			wantWriter: `[{"start":"2024-12-25T00:00:00Z","end":"2024-12-26T00:00:00Z","primary":"a","secondary":"b"},{"start":"2024-12-26T00:00:00Z","end":"2024-12-27T00:00:00Z","primary":"c","secondary":"d"}]`,
			wantErr:    false,
		},
		{
			name: "When plan is empty, throw an error",
			args: args{
				plan: []apis.Shift{},
			},
			wantWriter: "[]",
			wantErr:    false,
		},
		{
			name: "When plan is nil, do nothing",
			args: args{
				plan: nil,
			},
			wantWriter: "",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewJSONExporter()
			writer := &strings.Builder{}
			err := e.Write(tt.args.plan, writer)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
			}
			if content := writer.String(); content != tt.wantWriter {
				t.Errorf("Write() json = %v, want = %v", content, tt.wantWriter)
			}
		})
	}
}
