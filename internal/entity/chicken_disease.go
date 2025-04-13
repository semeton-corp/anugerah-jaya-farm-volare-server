package entity

type ChickenDiseaseMonitoring struct {
	Id                  uint64  `gorm:"primary_key;auto_increment"`
	ChickenMonitoringId uint64  `gorm:"type:bigint;not null"`
	Disease             string  `gorm:"type:varchar(255);not null"`
	Medicine            string  `gorm:"type:varchar(255);not null"`
	Dose                float64 `gorm:"type:decimal;not null"`
	Unit                string  `gorm:"type:varchar(255);not null"`
	CreatedAt           string  `gorm:"type:timestamp;auto_create_time"`
	UpdatedAt           string  `gorm:"type:timestamp;auto_update_time"`
}
