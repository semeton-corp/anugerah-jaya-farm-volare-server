package entity

type Loan struct {
	ID        uint64  `gorm:"primary_key;auto_increment"`
	StaffID   uint64  `gorm:"not null"`
	Amount    float64 `gorm:"not null"`
	Status    string  `gorm:"type:varchar(255);not null"`
	Type      string  `gorm:"type:varchar(255)"`
	CreatedAt string  `gorm:"type:timestamp;auto_create_time"`
	CreatedBy string  `gorm:"type:varchar(26);not null"`
}
