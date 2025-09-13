package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/shopspring/decimal"
)

type Expense struct {
	Id                  uint64               `gorm:"primaryKey;autoIncrement"`
	ExpenseCategory     enum.ExpenseCategory `gorm:"type:int;not null"`
	Name                string               `gorm:"varchar(255);not null"`
	ReceiverName        string               `gorm:"type:varchar(255);not null"`
	ReceiverPhoneNumber string               `gorm:"type:varchar(255);not null"`
	Nominal             decimal.Decimal      `gorm:"type:decimal;not null"`
	PaymentMethod       enum.PaymentMethod   `gorm:"type:int;not null"`
	PaymentProof        string               `gorm:"type:text;not null"`
	Description         string               `gorm:"type:text;not null"`
	LocationId          uint64               `gorm:"type:bigint;not null"`
	Location            Location             `gorm:"foreignKey:LocationId;references:Id;constraint:OnDelete:CASCADE"`
	WarehouseId         sql.NullInt64        `gorm:"type:bigint"`
	Warehouse           Warehouse            `gorm:"foreignKey:WarehouseId;references:Id;constraint:OnDelete:CASCADE"`
	StoreId             sql.NullInt64        `gorm:"type:bigint"`
	Store               Store                `gorm:"foreignKey:StoreId;references:Id;constraint:OnDelete:CASCADE"`
	CageId              sql.NullInt64        `gorm:"type:bigint"`
	Cage                Cage                 `gorm:"foreignKey:CageId;refereces:Id;constraint:OnDelete:CASCADE"`
	LocationType        enum.LocationType    `gorm:"int;not null"`
	CreatedAt           time.Time            `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy           uuid.NullUUID        `gorm:"type:varchar(255)"`
	UpdatedAt           time.Time            `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy           uuid.NullUUID        `gorm:"type:varchar(255)"`
	CreatedByUser       User                 `gorm:"foreignKey:CreatedBy;references:Id;constraint:OnDelete:SET NULL"`
}
