package export

import (
	"strings"
	"testing"
	"time"

	"github.com/orltom/on-call-schedule/pkg/apis"
)

func TestCVSExporter_Write(t *testing.T) {
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
			wantWriter: `Start,End,Primary,Secondary
2024-12-25,2024-12-26,a,b
2024-12-26,2024-12-27,c,d
`,
			wantErr: false,
		},
		{
			name: "When plan is empty, print CSV header",
			args: args{
				plan: []apis.Shift{},
			},
			wantWriter: "Start,End,Primary,Secondary\n",
			wantErr:    false,
		},
		{
			name: "When plan is nil, print CSV header",
			args: args{
				plan: nil,
			},
			wantWriter: "Start,End,Primary,Secondary\n",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewCVSCExporter()
			writer := &strings.Builder{}
			err := e.Write(tt.args.plan, writer)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
			}
			if content := writer.String(); content != tt.wantWriter {
				t.Errorf("Write() cvs = %v, want = %v", content, tt.wantWriter)
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
