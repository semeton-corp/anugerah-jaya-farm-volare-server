package dto

import (
	"time"

	"github.com/google/uuid"
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
	UserId         uuid.UUID
	Page           uint64                    `query:"page"`
	PresenceStatus param.PresenceStatusParam `query:"presenceStatus"`
	Month          param.MonthParam          `query:"month" validate:"required"`
	Year           uint64                    `query:"year" validate:"required"`
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
	TotalPage uint64                 `json:"totalPage,omitempty"`
	TotalData uint64                 `json:"totalData,omitempty"`
	Presences []PresenceListResponse `json:"presences"`
}

type RoleLocationPresenceSummaryResponse struct {
	RoleId              uint64 `json:"roleId"`
	RoleName            string `json:"roleName"`
	PlaceId             uint64 `json:"placeId"`
	PlaceName           string `json:"placeName"`
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

type GetUserPresenceSummaryFilter struct {
	RoleId    uint64                      `query:"roleId" validate:"roleId"`
	PlaceType param.LocationWorkTypeParam `query:"placeType" validate:"required"`
	PlaceId   uint64                      `query:"placeId" validate:"required"`
	Month     param.MonthParam            `query:"month" validate:"required"`
	Year      uint64                      `query:"year" validate:"required"`
}

type GetUserPresenceWorkDetailSummaryFilter struct {
	RoleId    uint64                      `query:"roleId" validate:"roleId"`
	PlaceType param.LocationWorkTypeParam `query:"placeType" validate:"required"`
	PlaceId   uint64                      `query:"placeId" validate:"required"`
	Date      param.DateParam             `query:"date"`
}

type UserPresenceSummaryResponse struct {
	UserId              string `json:"id"`
	UserName            string `json:"name"`
	UserPhotoProfile    string `json:"photoProfile"`
	UserEmail           string `json:"email"`
	RoleName            string `json:"roleName"`
	TotalPresentUser    uint64 `json:"totalPresentUser"`
	TotalSickUser       uint64 `json:"totalSickUser"`
	TotalPermissionUser uint64 `json:"totalPermissionUser"`
	TotalAlphaUser      uint64 `json:"totalAlphaUser"`
}

type UserPresenceWorkDetailSummaryResponse struct {
	UserId             string  `json:"id"`
	UserName           string  `json:"name"`
	UserPhotoProfile   string  `json:"photoProfile"`
	UserEmail          string  `json:"email"`
	RoleName           string  `json:"roleName"`
	PresenceStatus     string  `json:"status"`
	ArrivedTime        string  `json:"arrivedTime"`
	DepartureTime      string  `json:"departureTime"`
	WorkDonePercentage float64 `json:"WorkDonePercentage"`
}

type ApprovalPresenceRequest struct {
	AcceptedUserIds []string `json:"acceptedUserIds"`
	RejectedUserIds []string `json:"rejectedUserIds"`
}

type GetUserPresenceFilter struct {
	UserIds []uuid.UUID     `query:"userIds"`
	Date    param.DateParam `query:"date"`
}
