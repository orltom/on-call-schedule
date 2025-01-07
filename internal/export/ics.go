package export

import (
	_ "embed"
	"errors"
	"fmt"
	"io"
	"text/template"

	"github.com/orltom/on-call-schedule/pkg/apis"
)

//go:embed ics.tmpl
var templateICSFile string

var _ apis.Exporter = &TableExporter{}

type event struct {
	apis.Shift
	UID         string
	Category    string
	XColor      string
	Attendee    string
	Summary     string
	Description string
}

type ICSExporter struct {
	template *template.Template
}

func NewICSExporter() *ICSExporter {
	tmpl, _ := template.New("ical").Parse(templateICSFile)
	return &ICSExporter{template: tmpl}
}

func (e *ICSExporter) Write(plan []apis.Shift, writer io.Writer) error {
	if plan == nil {
		return errors.New("shifts must not be nil")
	}

	events := make([]event, 0, len(plan))
	for idx := range plan {
		s := plan[idx]
		events = append(events, mapPrimaryEvent(s), mapSecondaryEvent(s))
	}

	if err := e.template.ExecuteTemplate(writer, "icalCalendar", events); err != nil {
		return fmt.Errorf("can not generate ICS: %w", err)
	}
	return nil
}

func mapPrimaryEvent(shift apis.Shift) event {
	return event{
		Shift:       shift,
		UID:         "primary-" + shift.Start.Format("20060102"),
		Category:    "Primary",
		XColor:      "#FF5733",
		Attendee:    string(shift.Primary),
		Summary:     "Primary On-Call: " + string(shift.Primary),
		Description: fmt.Sprintf("Primary: %s\\nSecondary: %s", shift.Primary, shift.Secondary),
	}
}

func mapSecondaryEvent(shift apis.Shift) event {
	return event{
		Shift:       shift,
		UID:         "secondary-" + shift.Start.Format("20060102"),
		Category:    "Secondary",
		XColor:      "#33FF57",
		Attendee:    string(shift.Secondary),
		Summary:     "Secondary On-Call: " + string(shift.Secondary),
		Description: fmt.Sprintf("Primary: %s\\nSecondary: %s", shift.Primary, shift.Secondary),
	}
}
