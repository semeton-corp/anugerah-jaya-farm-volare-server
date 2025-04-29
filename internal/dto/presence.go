package dto

import "github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/param"

type PresenceListResponse struct {
	Id        uint64        `json:"id"`
	Staff     StaffResponse `json:"staff"`
	Date      string        `json:"date"`
	StartTime string        `json:"startTime"`
	EndTime   string        `json:"endTime"`
	ExtraTime string        `json:"extraTime"`
	IsPresent bool          `json:"isPresent"`
}

type PresenceResponse struct {
	Id        uint64        `json:"id"`
	Staff     StaffResponse `json:"staff"`
	Date      string        `json:"date"`
	StartTime string        `json:"startTime"`
	EndTime   string        `json:"endTime"`
	ExtraTime string        `json:"extraTime"`
	IsPresent bool          `json:"isPresent"`
}

type GetPresenceFilter struct {
	Status param.PresenceFilterParam `query:"status"`
}

type UpdateStaffPresenceRequest struct {
	IsPresent bool   `json:"isPresent" validate:"required"`
	StartTime string `json:"startTime"` // format: "15:04"`
	EndTime   string `json:"endTime"`
}
