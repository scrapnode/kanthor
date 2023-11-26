package utils

import (
	"time"

	"github.com/spf13/cobra"
)

var format = "2006-01-02"
var formatlong = format + " 15:04:05"

func DateRange(cmd *cobra.Command) (*time.Time, *time.Time, error) {
	f, err := cmd.Flags().GetString("from")
	if err != nil {
		return nil, nil, err
	}
	from, err := time.Parse(format, f)
	if err != nil {
		from, err = time.Parse(formatlong, f)
		if err != nil {
			return nil, nil, err
		}
	}

	t, err := cmd.Flags().GetString("to")
	if err != nil {
		return nil, nil, err
	}
	to, err := time.Parse(format, t)
	if err != nil {
		to, err = time.Parse(formatlong, f)
		if err != nil {
			return nil, nil, err
		}
	}

	return &from, &to, nil
}
