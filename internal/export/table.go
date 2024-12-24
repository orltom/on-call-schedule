package export

import (
	_ "embed"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/template"
	"time"

	"github.com/orltom/on-call-schedule/pkg/apis"
)

var ErrNotDate = errors.New("invalid date")

//go:embed table.tmpl
var templateTableFile string

var _ apis.Exporter = &TableExporter{}

type TableExporter struct {
	template *template.Template
}

func NewTableExporter() *TableExporter {
	fm := template.FuncMap{"dateOnly": formatDateOnly, "trunc": truncate}
	tmpl, _ := template.New("table-export").Funcs(fm).Parse(templateTableFile)
	return &TableExporter{template: tmpl}
}

func (e *TableExporter) Write(plan []apis.Shift, writer io.Writer) error {
	if err := e.template.Execute(writer, plan); err != nil {
		return fmt.Errorf("can not generate CVS: %w", err)
	}
	return nil
}

func truncate(item reflect.Value) (string, error) {
	var maxLen = 15
	v := fmt.Sprintf("%v", item.Interface())
	if len(v) == maxLen {
		return v, nil
	}

	if len(v) > maxLen {
		return v[:12] + "...", nil
	}

	return v + strings.Repeat(" ", maxLen-len(v)-1), nil
}

func formatDateOnly(item reflect.Value) (string, error) {
	if item.Kind() != reflect.Struct {
		return "", ErrNotDate
	}

	t, ok := item.Interface().(time.Time)
	if !ok {
		return "", ErrNotDate
	}

	return t.Format(time.DateOnly), nil
}
