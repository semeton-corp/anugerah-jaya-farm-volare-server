package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/param"
)

type PresenceListResponse struct {
	Id               uint64       `json:"id"`
	User             UserResponse `json:"user"`
	Date             string       `json:"date"`
	StartTime        string       `json:"startTime"`
	EndTime          string       `json:"endTime"`
	Overtime         float64      `json:"overTime"`
	Status           string       `json:"status"`
	SubmissionStatus string       `json:"submissionStatus"`
	Evidence         string       `json:"evidence"`
	Note             string       `json:"note"`
}

type PresenceResponse struct {
	Id               uint64       `json:"id"`
	User             UserResponse `json:"user"`
	Date             string       `json:"date"`
	StartTime        string       `json:"startTime"`
	EndTime          string       `json:"endTime"`
	Overtime         float64      `json:"overTime"`
	Status           string       `json:"status"`
	SubmissionStatus string       `json:"submissionStatus"`
	Evidence         string       `json:"evidence"`
	Note             string       `json:"note"`
	CreatedAt        time.Time    `json:"-"`
}

type GetPresenceFilter struct {
	UserId uuid.UUID
	Page   uint64           `query:"page"`
	Month  param.MonthParam `query:"month" validate:"required"`
	Year   uint64           `query:"year" validate:"required"`
}

type UpdateUserPresenceRequest struct {
	Status    string  `json:"status" validate:"required,presenceStatus"`
	StartTime string  `json:"startTime"` // format: "15:04"`
	EndTime   string  `json:"endTime"`
	Evidence  string  `json:"evidence"`
	Note      string  `json:"note"`
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
}

type PresenceListPaginationResponse struct {
	TotalPage uint64                 `json:"totalPage"`
	TotalData uint64                 `json:"totalData"`
	Presences []PresenceListResponse `json:"presences"`
}

type LocationPresenceSummaryResponse struct {
	Place               string `json:"place"`
	PlaceType           string `json:"placeType"`
	TotalUser           uint64 `json:"totalUser"`
	TotalPresentUser    uint64 `json:"totalPresentUser"`
	TotalSickUser       uint64 `json:"totalSickUser"`
	TotalPermissionUser uint64 `json:"totalPermissionUser"`
	TotalAlphaUser      uint64 `json:"totalAlphaUser"`
}

type GetLocationPresenceSummaryFilter struct {
	Date param.DateParam `query:"date"`
}

type PresenceSummaryDetailResponse struct {
	Month enum.Month `query:"month"`
	Year  uint64     `query:"year"`
}

type GetPresenceSummaryDetailFilter struct {
	PlaceType string `query:"placeType" validate:"required"`
	PlaceId   uint64 `query:"placeId" validate:"required"`
}
