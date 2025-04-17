package param

import "time"

type DateParam time.Time

func (cd *DateParam) UnmarshalText(text []byte) error {
	parsedTime, err := time.Parse("2006-01-02", string(text))
	if err != nil {
		return err
	}
	*cd = DateParam(parsedTime)
	return nil
}

func (cd DateParam) Value() time.Time {
	return time.Time(cd)
}
