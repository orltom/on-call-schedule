package export

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/orltom/on-call-schedule/pkg/apis"
)

var _ apis.Exporter = &JSONExporter{}

type JSONExporter struct{}

func NewJSONExporter() *JSONExporter {
	return &JSONExporter{}
}

func (c *JSONExporter) Write(plan []apis.Shift, writer io.Writer) error {
	if plan == nil {
		return errors.New("plan must not be nil")
	}

	content, err := json.Marshal(plan)
	if err != nil {
		return fmt.Errorf("could not marshal plan: %w", err)
	}

	_, err = writer.Write(content)
	if err != nil {
		return fmt.Errorf("could not write plan: %w", err)
	}

	return nil
}
