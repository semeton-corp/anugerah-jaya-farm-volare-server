package dto

import (
	"time"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/param"
)

type WorkUserResponse struct {
	DailyWorks      []DailyWorkUserResponse      `json:"dailyWorks"`
	AdditionalWorks []AdditionalWorkUserResponse `json:"additionalWorks"`
}

type CreateDailyWorkRequest struct {
	RoleId          uint64                         `json:"roleId" validate:"required"`
	DailyWorkDetail []CreateDailyWorkDetailRequest `json:"dailyWorkDetail"`
}

type CreateDailyWorkDetailRequest struct {
	Id          uint64 `json:"id"`
	Description string `json:"description" validate:"required"`
	StartTime   string `json:"startTime" validate:"required"`
	EndTime     string `json:"endTime" validate:"required"`
}

type DailyWorkResponse struct {
	Role       RoleResponse              `json:"role"`
	DailyWorks []DailyWorkDetailResponse `json:"dailyWorks"`
}

type DailyWorkDetailResponse struct {
	Id          uint64 `json:"id"`
	Description string `json:"description"`
	Date        string `json:"date"`
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
}

type CreateAdditionalWorkRequest struct {
	Name         string   `json:"name" validate:"required"`
	LocationId   uint64   `json:"locationId" validate:"required"`
	LocationType string   `json:"locationType" validate:"required,locationType"`
	PlaceId      *uint64  `json:"placeId"`
	Description  string   `json:"description" validate:"required"`
	Slot         uint64   `json:"slot" validate:"required"`
	Salary       string   `json:"salary" validate:"required"`
	WorkDate     string   `json:"workDate" validate:"required"`
	UserIds      []string `json:"userIds"`
}

type UpdateAdditionalWorkRequest struct {
	Name         string   `json:"name" validate:"required"`
	LocationId   uint64   `json:"locationId" validate:"required"`
	LocationType string   `json:"locationType" validate:"required,locationType"`
	PlaceId      *uint64  `json:"placeId"`
	Description  string   `json:"description" validate:"required"`
	Slot         uint64   `json:"slot" validate:"required"`
	Salary       string   `json:"salary" validate:"required"`
	WorkDate     string   `json:"workDate" validate:"required"`
	UserIds      []string `json:"userIds"`
}

type AdditionalWorkResponse struct {
	Id                            uint64                                  `json:"id"`
	Name                          string                                  `json:"name"`
	Location                      LocationResponse                        `json:"location"`
	LocationType                  string                                  `json:"locationType"`
	Description                   string                                  `json:"description"`
	Place                         string                                  `json:"place"`
	Date                          string                                  `json:"date"`
	Time                          string                                  `json:"time"`
	Slot                          uint64                                  `json:"slot"`
	Salary                        string                                  `json:"salary"`
	AdditionalWorkUserInformation []AdditionalWorkUserInformationResponse `json:"additionalWorkUserInformation"`
}

type AdditionalWorkUserInformationResponse struct {
	UserId   string `json:"userId"`
	RoleId   uint64 `json:"roleId"`
	RoleName string `json:"roleName"`
	UserName string `json:"userName"`
}

type AdditionalWorkDetailResponse struct {
	Id          uint64           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Date        string           `json:"date"`
	Time        string           `json:"time"`
	Salary      string           `json:"salary"`
	Location    LocationResponse `json:"location"`
	PlaceType   string           `json:"placeType"`
	Place       string           `json:"place"`
}

type AdditionalWorkUserResponse struct {
	Id             uint64                       `json:"id"`
	IsDone         bool                         `json:"isDone"`
	Note           string                       `json:"note"`
	AdditionalWork AdditionalWorkDetailResponse `json:"additionalWork"`
	TakenDate      string                       `json:"finishedAt"`
	TakenTime      string                       `json:"finishedTime"`
	CreatedAt      time.Time                    `json:"-"`
}

type AdditionalWorkUserListPaginationResponse struct {
	TotalPage           uint64                       `json:"totalPage,omitempty"`
	TotalData           uint64                       `json:"totalData,omitempty"`
	AdditionalWorkUsers []AdditionalWorkUserResponse `json:"additionalWorkUsers"`
}

type DailyWorkUserResponse struct {
	Id           uint64                  `json:"id"`
	IsDone       bool                    `json:"isDone"`
	Note         string                  `json:"note"`
	DailyWork    DailyWorkDetailResponse `json:"dailyWork"`
	FinishedDate string                  `json:"finishedDate"`
	FinishedTime string                  `json:"finishedTime"`
	CreatedAt    time.Time               `json:"-"`
	Status       string                  `json:"status"`
}

type DailyWorkUserListPaginationResponse struct {
	TotalPage      uint64                  `json:"totalPage,omitempty"`
	TotalData      uint64                  `json:"totalData,omitempty"`
	DailyWorkUsers []DailyWorkUserResponse `json:"dailyWorkUsers"`
}

type DailyWorkListResponse struct {
	Role      RoleResponse `json:"role"`
	TotalWork uint64       `json:"totalWork"`
	TotalUser uint64       `json:"totalUser"`
}

type AdditionalWorkListResponse struct {
	Id                   uint64 `json:"id"`
	Date                 string `json:"date"`
	Time                 string `json:"time"`
	Name                 string `json:"name"`
	Description          string `json:"description"`
	Location             string `json:"location"`
	Place                string `json:"place"`
	RemainingSlot        uint64 `json:"remainingSlot"`
	Status               string `json:"status"`
	IsTakenByCurrentUser *bool  `json:"IsTakenByCurrentUser,omitempty"`
}

type GetAdditonalWorkFilter struct {
	Status       string                  `query:"status"`
	LocationId   uint64                  `query:"locationId"`
	LocationType param.LocationTypeParam `query:"locationType"`
	PlaceIds     []uint64                `query:"placeIds"`
}

type UpdateAdditionalWorkUserRequest struct {
	IsDone bool   `json:"isDone"`
	Note   string `json:"note"`
}

type UpdateDailyWorkUserRequest struct {
	IsDone bool   `json:"isDone"`
	Note   string `json:"note"`
}

type GetDailyWorkUserFilter struct {
	Date        param.DateParam  `query:"date"`
	Month       param.MonthParam `query:"month"`
	Year        uint64           `query:"year"`
	Page        uint64           `query:"page"`
	WithDeleted *bool            `query:"withDeleted" validate:"required"`
}

type GetAdditionalWorkUserFilter struct {
	Month                param.MonthParam `query:"month"`
	Year                 uint64           `query:"year"`
	Page                 uint64           `query:"page"`
	IsAdditionalWorkFull bool             `query:"isAdditionalWorkFull"`
	WithDeleted          *bool            `query:"withDeleted" validate:"required"`
}

type WorkOveriew struct {
	AdditionalWorkSummaries []AdditionalWorkListResponse `json:"additionalWorkSummaries"`
	DailyWorkSummaries      []DailyWorkListResponse      `json:"dailyWorkSummaries"`
}
