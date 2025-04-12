package entity

type RequestProductStorage struct {
	ID        uint64 `gorm:"primary_key;auto_increment"`
	StorageID uint64 `gorm:"not null"`
	ProductID uint64 `gorm:"not null"`
	CreatedAt string `gorm:"type:timestamp;auto_create_time"`
	UpdatedAt string `gorm:"type:timestamp;auto_update_time"`
}
