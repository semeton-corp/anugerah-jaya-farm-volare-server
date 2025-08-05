package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/shopspring/decimal"
)

type WarehouseItemCornProcurementPayment struct {
	Id                             uint64             `gorm:"primaryKey;autoIncrement"`
	WarehouseItemCornProcurementId uint64             `gorm:"type:bigint;not null"`
	PaymentDate                    time.Time          `gorm:"type:date;not null"`
	Nominal                        decimal.Decimal    `gorm:"type:decimal;not null"`
	PaymentProof                   string             `gorm:"type:text;not null"`
	PaymentMethod                  enum.PaymentMethod `gorm:"type:int;not null"`
	CreatedAt                      time.Time          `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy                      uuid.NullUUID      `gorm:"type:varchar(255)"`
	UpdatedAt                      time.Time          `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy                      uuid.NullUUID      `gorm:"type:varchar(255)"`
}
