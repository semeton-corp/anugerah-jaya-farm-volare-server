package entity

type SellIncome struct {
	ID                  uint64 `gorm:"primaryKey;autoIncrement"`
	Location            string `gorm:"type:varchar(255);not null"`
	Type                string `gorm:"type:varchar(255);not null"`
	CustomerName        string `gorm:"type:varchar(255);not null"`
	CustomerPhoneNumber string `gorm:"type:varchar(15);not null"`
	CustomerAddress     string `gorm:"type:text;not null"`
	CreatedAt           string `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy           string `gorm:"type:varchar(26);not null"`
	UpdatedAt           string `gorm:"type:timestamp;autoUpdateTime"`
}
