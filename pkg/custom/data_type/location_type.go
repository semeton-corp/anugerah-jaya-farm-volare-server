package datatype

import (
	"database/sql"
	"database/sql/driver"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
)

type NullLocationType struct {
	LocationType enum.LocationType
	Valid        bool
}

func (n *NullLocationType) Scan(value interface{}) error {
	if value == nil {
		n.LocationType, n.Valid = enum.LocationTypeUnknown, false
		return nil
	}

	var i sql.NullInt64
	if err := i.Scan(value); err != nil {
		return err
	}

	if !i.Valid {
		n.LocationType, n.Valid = enum.LocationTypeUnknown, false
		return nil
	}

	n.LocationType = enum.LocationType(i.Int64)
	n.Valid = true
	return nil
}

func (n NullLocationType) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return int64(n.LocationType), nil
}

func (n NullLocationType) String() string {
	if !n.Valid {
		return "NULL"
	}
	return n.LocationType.String()
}
