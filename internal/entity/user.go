package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/shopspring/decimal"
)

type User struct {
	Id             uuid.UUID           `gorm:"type:varchar(255);primaryKey"`
	Username       string              `gorm:"type:varchar(255);not null;unique"`
	Email          string              `gorm:"type:varchar(255);not null;unique"`
	Password       string              `gorm:"type:varchar(255);not null"`
	LocationId     sql.NullInt64       `gorm:"type:bigint"`
	Location       Location            `gorm:"foreignKey:LocationId;references:Id;constraint:OnDelete:CASCADE"`
	RoleId         uint64              `gorm:"type:bigint;not null"`
	Role           Role                `gorm:"foreignKey:RoleId;references:Id;constraint:OnDelete:set null"`
	PhotoProfile   string              `gorm:"type:text;default:null"`
	Name           string              `gorm:"type:varchar(255);not null"`
	PhoneNumber    string              `gorm:"type:varchar(255);not null"`
	Address        string              `gorm:"type:text;not null"`
	SalaryInterval enum.SalaryInterval `gorm:"type:int;not null;default:1"`
	Salary         decimal.Decimal     `gorm:"type:decimal;not null"`
	CreatedByOwner uuid.NullUUID       `gorm:"type:varchar(255)"`
	CreatedAt      time.Time           `gorm:"type:timestamp;autoCreateTime"`
	UpdatedBy      uuid.NullUUID       `gorm:"type:varchar(255)"`
	UpdatedAt      time.Time           `gorm:"type:timestamp;autoUpdateTime"`
}
