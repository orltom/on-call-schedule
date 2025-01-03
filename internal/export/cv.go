package export

import (
	"encoding/csv"
	"fmt"
	"io"
	"time"

	"github.com/orltom/on-call-schedule/pkg/apis"
)

var _ apis.Exporter = &CVSExporter{}

type CVSExporter struct{}

func NewCVSCExporter() *CVSExporter {
	return &CVSExporter{}
}

func (c *CVSExporter) Write(shifts []apis.Shift, writer io.Writer) error {
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	if err := csvWriter.Write([]string{"Start", "End", "Primary", "Secondary"}); err != nil {
		return fmt.Errorf("write csv header: %w", err)
	}

	for i := range shifts {
		s := shifts[i]
		record := []string{s.Start.Format(time.DateOnly), s.End.Format(time.DateOnly), string(s.Primary), string(s.Secondary)}
		if err := csvWriter.Write(record); err != nil {
			return fmt.Errorf("write csv data: %w", err)
		}
	}

	return nil
}
