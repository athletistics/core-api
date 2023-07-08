// Brief   : Date only data adaptor for GORM
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

// DateFormat is the format template according to https://pkg.go.dev/time#pkg-constants
// The time package in Go 1.20 defines this constant (DateOnly).
const DateFormat = "2006-01-02"
const NdsDateFormat = "02.01.2006"

type Date time.Time

func NewDate(year int, month time.Month, day int) Date {
	t := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return Date(t)
}

func (d Date) T() time.Time {
	return time.Time(d)
}

func (d Date) ForNds() string {
	return time.Time(d).Format(NdsDateFormat)
}

func (d *Date) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		return d.UnmarshalText(string(v))
	case string:
		return d.UnmarshalText(v)
	case time.Time:
		*d = Date(v)
	case nil:
		*d = Date{}
	default:
		return fmt.Errorf("cannot sql.Scan() MyTime from: %#v", v)
	}
	return nil
}

func (d Date) Value() (driver.Value, error) {
	return driver.Value(time.Time(d).Format(DateFormat)), nil
}

func (d *Date) UnmarshalText(value string) error {
	parsedDate, err := time.Parse(DateFormat, value)
	if err != nil {
		return err
	}
	*d = Date(parsedDate)
	return nil
}

func (Date) GormDataType() string {
	return "DATE"
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%v"`, time.Time(d).Format(DateFormat))), nil

}

func (d *Date) UnmarshalJSON(b []byte) error {
	if string(b) == "null" || string(b) == `""` {
		return nil
	}

	var dateString string
	if err := json.Unmarshal(b, &dateString); err != nil {
		return err
	}

	return d.UnmarshalText(dateString)
}
