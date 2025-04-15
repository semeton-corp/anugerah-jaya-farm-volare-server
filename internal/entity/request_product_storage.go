package entity

type RequestProductStorage struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement"`
	StorageID uint64 `gorm:"not null"`
	ProductID uint64 `gorm:"not null"`
	CreatedAt string `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt string `gorm:"type:timestamp;autoUpdateTime"`
}
