package datatype

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type TimeOnly struct {
	Time *time.Time
}

func (TimeOnly) GormDataType() string {
	return "time"
}

func (TimeOnly) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return "time"
}

func (timeOnly TimeOnly) Value() (driver.Value, error) {
	if timeOnly.Time != nil && !timeOnly.Time.IsZero() {
		return timeOnly.Time.Format("15:04"), nil
	}
	return nil, nil
}

func (timeOnly *TimeOnly) GetTime() *time.Time {
	return timeOnly.Time
}

func (timeOnly *TimeOnly) Scan(value interface{}) error {
	if value == nil {
		timeOnly.Time = nil
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		t := time.Date(0, 1, 1, v.Hour(), v.Minute(), v.Second(), v.Nanosecond(), time.UTC)
		timeOnly.Time = &t
		return nil
	case []byte:
		scannedString := string(v)
		return timeOnly.parseString(scannedString)
	case string:
		return timeOnly.parseString(v)
	default:
		return errx.InternalServerError("failed to scan time: unknown type")
	}
}

func (timeOnly *TimeOnly) parseString(s string) error {
	if s == "" {
		timeOnly.Time = nil
		return nil
	}
	parsedTime, err := time.Parse("15:04:05", s)
	if err != nil {
		parsedTime, err = time.Parse("15:04:05Z07:00", s)
		if err != nil {
			return errx.InternalServerError("failed to parse time string: " + err.Error())
		}
	}
	timeOnly.Time = &parsedTime
	return nil
}

func (timeOnly TimeOnly) MarshalJSON() ([]byte, error) {
	if timeOnly.Time == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(timeOnly.Time.Format("15:04:05"))
}

func (timeOnly *TimeOnly) UnmarshalJSON(bs []byte) error {
	var s *string
	err := json.Unmarshal(bs, &s)
	if err != nil {
		return err
	}
	if s == nil || *s == "" {
		timeOnly.Time = nil
		return nil
	}
	t, err := time.ParseInLocation("15:04", *s, time.UTC)
	if err != nil {
		return err
	}
	timeOnly.Time = &t
	return nil
}

func ParseTimeOnly(value string) (TimeOnly, error) {
	if value == "" {
		return TimeOnly{nil}, nil
	}
	t, err := time.Parse("15:04", value)
	if err != nil {
		return TimeOnly{}, err
	}
	return TimeOnly{&t}, nil
}
