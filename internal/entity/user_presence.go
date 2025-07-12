package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	datatype "github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/custom/data_type"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
)

type UserPresence struct {
	Id                       uint64                        `gorm:"primaryKey;autoIncrement"`
	UserId                   uuid.UUID                     `gorm:"type:bigint;not null"`
	User                     User                          `gorm:"foreignKey:UserId;references:Id;constraint:OnDelete:CASCADE"`
	StartTime                datatype.TimeOnly             `gorm:"type:timestamp"`
	EndTime                  datatype.TimeOnly             `gorm:"type:timestamp"`
	Status                   enum.PresenceStatus           `gorm:"type:int;not null"`
	Note                     sql.NullString                `gorm:"type:text"`
	Evidence                 sql.NullString                `gorm:"type:text"`
	SubmissionPresenceStatus enum.SubmissionPresenceStatus `gorm:"int;not null"`
	CreatedAt                time.Time                     `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy                uuid.NullUUID                 `gorm:"type:varchar(255)"`
	UpdatedAt                time.Time                     `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy                uuid.NullUUID                 `gorm:"type:varchar(255)"`
}

type LocationPresenceSummary struct {
	PlaceId        uint64              `gorm:"column:place_id"`
	PlaceName      string              `gorm:"column:place_name"`
	UserId         uuid.UUID           `gorm:"column:user_id"`
	PresenceStatus enum.PresenceStatus `gorm:"column:presence_status"`
}

type UserPresenceSummary struct {
	ID            uint64 `gorm:"column:id"`
	Name          string `gorm:"column:name"`
	PhotoProfile  string `gorm:"column:photo_profile"`
	Email         string `gorm:"column:email"`
	RoleName      string `gorm:"column:role_name"`
	Status        string `gorm:"column:status"`
	TotalPresence int64  `gorm:"column:total_presence"`
}
