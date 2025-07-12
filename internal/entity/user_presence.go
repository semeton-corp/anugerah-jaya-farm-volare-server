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
	UserId           uuid.UUID `gorm:"column:user_id"`
	UserName         string    `gorm:"column:user_name"`
	UserPhotoProfile string    `gorm:"column:user_photo_profile"`
	UserEmail        string    `gorm:"column:user_email"`
	RoleName         string    `gorm:"column:role_name"`
	TotalPresent     int64     `gorm:"column:total_present"`
	TotalSick        int64     `gorm:"column:total_sick"`
	TotalPermission  int64     `gorm:"column:total_permission"`
	TotalAlpha       int64     `gorm:"column:total_alpha"`
}
