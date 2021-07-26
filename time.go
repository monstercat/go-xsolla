package xsolla

import (
	"time"
)

var (
	// Depending on the route sometimes the timezone offset has colon and sometimes it
	// does not.
	timezoneLayouts = []string{
		`"2006-01-02T15:04:05-07:00"`,
		`"2006-01-02T15:04:05-0700"`,
	}
)

type Time time.Time

func (t *Time) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	var tim time.Time
	var err error
	for _, layout := range timezoneLayouts {
		tim, err = time.Parse(layout, string(data))
		if err == nil {
			break
		}
	}
	if err != nil {
		return err
	}
	*t = Time(tim)
	return nil
}
