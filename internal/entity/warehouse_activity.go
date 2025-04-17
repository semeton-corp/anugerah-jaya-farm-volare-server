package entity

type WarehouseActivity struct {
	Id          uint64 `gorm:"primaryKey;autoIncrement"`
	Description string `gorm:"type:varchar(255);not null"`
	Status      uint8  `gorm:"type:int;not null"`
	CreatedAt   string `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy   string `gorm:"type:varchar(255);not null"`
}
