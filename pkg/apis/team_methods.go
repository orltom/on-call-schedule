package apis

import (
	"strings"
	"time"
)

func (v *VacationDay) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		return nil
	}

	date := strings.Trim(string(data), "\"")
	t, err := time.Parse(time.DateOnly, date)
	if err != nil {
		return err
	}
	*v = VacationDay{t}
	return nil
}
