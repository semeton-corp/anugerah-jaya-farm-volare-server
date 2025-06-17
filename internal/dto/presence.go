package dto

import (
	"time"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/param"
)

type PresenceListResponse struct {
	Id        uint64       `json:"id"`
	User      UserResponse `json:"user"`
	Date      string       `json:"date"`
	StartTime string       `json:"startTime"`
	EndTime   string       `json:"endTime"`
	Overtime  float64      `json:"overTime"`
	IsPresent bool         `json:"isPresent"`
}

type PresenceResponse struct {
	Id        uint64       `json:"id"`
	User      UserResponse `json:"user"`
	Date      string       `json:"date"`
	StartTime string       `json:"startTime"`
	EndTime   string       `json:"endTime"`
	Overtime  float64      `json:"overTime"`
	IsPresent bool         `json:"isPresent"`
	CreatedAt time.Time    `json:"-"`
}

type GetPresenceFilter struct {
	Page  uint64           `query:"page"`
	Month param.MonthParam `query:"month" validate:"required"`
	Year  uint64           `query:"year" validate:"required"`
}

type UpdateStaffPresenceRequest struct {
	IsPresent bool   `json:"isPresent" validate:"required"`
	StartTime string `json:"startTime"` // format: "15:04"`
	EndTime   string `json:"endTime"`
}

type PresenceListPaginationResponse struct {
	TotalPage uint64                 `json:"totalPage"`
	TotalData uint64                 `json:"totalData"`
	Presences []PresenceListResponse `json:"presences"`
}
