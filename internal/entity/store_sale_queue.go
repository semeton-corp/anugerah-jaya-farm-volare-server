package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
)

type StoreSaleQueue struct {
	Id                  uint64            `gorm:"primaryKey;autoIncrement"`
	CustomerId          sql.NullInt64     `gorm:"type:bigint"`
	Customer            Customer          `gorm:"foreignKey:CustomerId;references:Id;constraint:OnDelete:CASCADE"`
	CustomerName        sql.NullString    `gorm:"type:varchar(255)"`
	CustomerPhoneNumber sql.NullString    `gorm:"type:varchar(255)"`
	CustomerType        enum.CustomerType `gorm:"type:int;not null"`
	ItemId              uint64            `gorm:"type:bigint;not null"`
	Item                Item              `gorm:"foreignKey:ItemId;references:Id;constraint:OnDelete:CASCADE"`
	StoreId             uint64            `gorm:"type:bigint;not null"`
	Store               Store             `gorm:"foreignKey:StoreId;references:Id;constraint:OnDelete:CASCADE"`
	SaleUnit            enum.SaleUnit     `gorm:"type:int;not null"`
	Quantity            float64           `gorm:"type:decimal;not null"`
	CreatedAt           time.Time         `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy           uuid.NullUUID     `gorm:"type:varchar(255)"`
	UpdatedAt           time.Time         `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy           uuid.NullUUID     `gorm:"type:varchar(255)"`
}
