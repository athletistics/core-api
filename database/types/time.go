// Brief   : Time only data adaptor for GORM
// Source  : https://stackoverflow.com/a/65459041 (Ezequiel Muns)
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 06.07.2023

package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// TimeFormat is the format template according to https://pkg.go.dev/time#pkg-constants
// The time package in Go 1.20 defines a similar constant (TimeOnly).
const TimeFormat = "15:04"
const DbTimeFormat = "15:04:05"

type Time time.Time

func NewTime(hour, min, sec int) Time {
	t := time.Date(0, time.January, 1, hour, min, sec, 0, time.UTC)
	return Time(t)
}

func (t Time) T() time.Time {
	return time.Time(t)
}

func (t Time) ForNds() string {
	return time.Time(t).Format(TimeFormat)
}

func (t *Time) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		return t.UnmarshalText(string(v), DbTimeFormat)
	case string:
		return t.UnmarshalText(v, DbTimeFormat)
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
	return driver.Value(time.Time(t).Format(DbTimeFormat)), nil
}

func (t *Time) UnmarshalText(value string, format string) error {
	paprsedTime, err := time.Parse(format, value)
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
	return []byte(fmt.Sprintf(`"%v"`, time.Time(t).Format(TimeFormat))), nil
}

func (t *Time) UnmarshalJSON(b []byte) error {
	if string(b) == "null" || string(b) == `""` {
		return nil
	}

	var timeString string
	if err := json.Unmarshal(b, &timeString); err != nil {
		return err
	}

	return t.UnmarshalText(timeString, TimeFormat)
}
