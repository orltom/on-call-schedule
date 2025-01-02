package cli

import (
	"errors"
	"flag"
	"fmt"
	"maps"
	"os"
	"strings"
	"time"
)

func RequiredFlagPassed(f *flag.FlagSet, names ...string) (bool, []string) {
	var missedFlags []string
	visited := make(map[string]bool)

	f.Visit(func(fl *flag.Flag) {
		visited[fl.Name] = true
	})

	for _, name := range names {
		if !visited[name] {
			missedFlags = append(missedFlags, name)
		}
	}

	return len(missedFlags) == 0, missedFlags
}

func TimeValueVar(t *time.Time) func(s string) error {
	return func(s string) error {
		return validateNonEmpty(s, func(s string) error {
			tmp, err := time.Parse(time.DateTime, s)
			if err != nil {
				return fmt.Errorf("invalid date time format %s", s)
			}
			*t = tmp

			return nil
		})
	}
}

func FilePathVar(path *string) func(s string) error {
	return func(s string) error {
		return validateNonEmpty(s, func(s string) error {
			info, err := os.Stat(s)
			if errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("%s not found", s)
			}
			if info.IsDir() {
				return fmt.Errorf("%s %s", s, "is not a file")
			}
			*path = s

			return nil
		})
	}
}

func EnumValueVar[T comparable](enums map[string]T, transform *T) func(s string) error {
	return func(s string) error {
		return validateNonEmpty(s, func(s string) error {
			for k, v := range enums {
				if strings.EqualFold(k, s) {
					*transform = v
					return nil
				}
			}
			var keys []string
			for k := range maps.Keys(enums) {
				keys = append(keys, k)
			}
			return fmt.Errorf("expected one of (%s)", strings.Join(keys, ", "))
		})
	}
}

func validateNonEmpty(s string, validator func(string) error) error {
	if len(strings.TrimSpace(s)) < 1 {
		return errors.New("value not set")
	}
	return validator(s)
}
