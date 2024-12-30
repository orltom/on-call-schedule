package export

import (
	"strings"
	"testing"

	"github.com/orltom/on-call-schedule/pkg/apis"
)

func TestTableExporter_Write(t *testing.T) {
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
			name: "Print table result according shift plan",
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
			wantWriter: `From       | To         | Primary        | Secondary
----------------------------------------------------------
2024-12-25 | 2024-12-26 | a              | b             
2024-12-26 | 2024-12-27 | c              | d             

`,
			wantErr: false,
		},
		{
			name: "When plan is empty, print table header",
			args: args{
				plan: []apis.Shift{},
			},
			wantWriter: "From       | To         | Primary        | Secondary\n----------------------------------------------------------\n\n",
			wantErr:    false,
		},
		{
			name: "When plan is nil, print table header",
			args: args{
				plan: nil,
			},
			wantWriter: "From       | To         | Primary        | Secondary\n----------------------------------------------------------\n\n",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewTableExporter()
			writer := &strings.Builder{}
			err := e.Write(tt.args.plan, writer)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
			}
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("Write() gotWriter = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}
