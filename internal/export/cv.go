package export

import (
	"fmt"
	"io"
	"time"

	"github.com/orltom/on-call-schedule/pkg/apis"
)

var _ apis.Exporter = &CVSExporter{}

type CVSExporter struct{}

func (c *CVSExporter) Write(shifts []apis.Shift, writer io.Writer) error {
	if _, err := fmt.Fprint(writer, "from,to,primary,secondary\n"); err != nil {
		return fmt.Errorf("write csv header: %w", err)
	}
	for i := range shifts {
		s := shifts[i]
		if _, err := fmt.Fprintf(writer, "\"%s\",\"%s\",\"%s\",\"%s\"\n", s.Start.Format(time.DateOnly), s.End.Format(time.DateOnly), s.Primary, s.Secondary); err != nil {
			return fmt.Errorf("write csv data: %w", err)
		}
	}
	return nil
}

func NewCVSCExporter() *CVSExporter {
	return &CVSExporter{}
}
