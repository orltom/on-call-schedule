package export

import (
	"bytes"
	"testing"

	"github.com/orltom/on-call-schedule/pkg/apis"
)

func TestICSExporter_Write(t *testing.T) {
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
			name: "When plan is nil, throw an error",
			args: args{
				plan: nil,
			},
			wantWriter: "",
			wantErr:    true,
		},
		{
			name: "Print ICS result according shift plan",
			args: args{
				plan: []apis.Shift{
					{
						Start:     date("2024-12-25"),
						End:       date("2024-12-26"),
						Primary:   "a",
						Secondary: "b",
					},
				},
			},
			wantWriter: `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//ocsctl//on-call-shift//EN
CALSCALE:GREGORIAN

BEGIN:VEVENT
UID:primary-20241225
DTSTART:20241225T000000Z
DTEND:20241226T000000Z
SUMMARY:Primary On-Call: a
DESCRIPTION:Primary: a\nSecondary: b
CATEGORIES:Primary
ATTENDEE:mailto:a
X-COLOR:#FF5733
BEGIN:VALARM
TRIGGER:-PT8H
ACTION:DISPLAY
DESCRIPTION:Reminder: Primary On-Call: a starts in 8 hours.
END:VALARM
END:VEVENT

BEGIN:VEVENT
UID:secondary-20241225
DTSTART:20241225T000000Z
DTEND:20241226T000000Z
SUMMARY:Secondary On-Call: b
DESCRIPTION:Primary: a\nSecondary: b
CATEGORIES:Secondary
ATTENDEE:mailto:b
X-COLOR:#33FF57
BEGIN:VALARM
TRIGGER:-PT8H
ACTION:DISPLAY
DESCRIPTION:Reminder: Secondary On-Call: b starts in 8 hours.
END:VALARM
END:VEVENT

END:VCALENDAR`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewICSExporter()
			writer := &bytes.Buffer{}
			err := e.Write(tt.args.plan, writer)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("Write() gotWriter = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}
