// Brief   : Time only data adaptor for GORM
// Source  : https://stackoverflow.com/a/65459041 (Ezequiel Muns)
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 06.07.2023

package types

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// TimeFormat is the format template according to https://pkg.go.dev/time#pkg-constants
// The time package in Go 1.20 defines this constant (TimeOnly).
const TimeFormat = "15:04:05"
const NdsTimeFormat = "15:04"

type Time time.Time

func NewTime(hour, min, sec int) Time {
	t := time.Date(0, time.January, 1, hour, min, sec, 0, time.UTC)
	return Time(t)
}

func (t Time) T() time.Time {
	return time.Time(t)
}

func (t Time) ForNds() string {
	return time.Time(t).Format(NdsTimeFormat)
}

func (t *Time) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		return t.UnmarshalText(string(v))
	case string:
		return t.UnmarshalText(v)
	case time.Time:
		*t = Time(v)
	case nil:
		*t = Time{}
	default:
		return fmt.Errorf("cannot sql.Scan() MyTime from: %#v", v)
	}
	return nil
}

func (t Time) Value() (driver.Value, error) {
	return driver.Value(time.Time(t).Format(TimeFormat)), nil
}

func (t *Time) UnmarshalText(value string) error {
	paprsedTime, err := time.Parse(TimeFormat, value)
	if err != nil {
		return err
	}
	*t = Time(paprsedTime)
	return nil
}

func (Time) GormDataType() string {
	return "TIME"
}

func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(t).Format(TimeFormat)), nil
}

func (t *Time) UnmarshalJSON(b []byte) error {
	return t.UnmarshalText(string(b))
}
