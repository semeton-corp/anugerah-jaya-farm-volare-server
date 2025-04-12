package entity

type Product struct {
	Id           uint64 `gorm:"primary_key;auto_increment"`
	StorageId    uint64 `gorm:"not null"`
	Name         string `gorm:"type:varchar(255);not null"`
	Description  string `gorm:"type:text;not null"`
	Type         string `gorm:"type:varchar(255);not null"` // bahan baku atau bahan jadi
	Stock        uint64 `gorm:"not null"`
	MinimumStock uint64 `gorm:"not null"`
	MaximalStock uint64 `gorm:"not null"`
	ExpiredDate  string `gorm:"type:date"`
	LeadDate     string `gorm:"type:date"`
	CreatedAt    string `gorm:"type:timestamp;auto_create_time"`
	UpdatedAt    string `gorm:"type:timestamp;auto_update_time"`
}
