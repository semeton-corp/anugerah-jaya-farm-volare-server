package dto

type WorkStaffResponse struct {
	DailyWorks      []DailyWorkStaffResponse      `json:"dailyWorks"`
	AdditionalWorks []AdditionalWorkStaffResponse `json:"additionalWorks"`
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
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
}

type CreateAdditionalWorkRequest struct {
	Description string `json:"description" validate:"required"`
	Location    string `json:"location" validate:"required"`
	Slot        uint64 `json:"slot" validate:"required"`
}

type UpdateAdditionalWorkRequest struct {
	Description string `json:"description" validate:"required"`
	Location    string `json:"location" validate:"required"`
	Slot        uint64 `json:"slot" validate:"required"`
}

type AdditionalWorkResponse struct {
	Id                             uint64                                   `json:"id"`
	Description                    string                                   `json:"description" validate:"required"`
	Location                       string                                   `json:"location" validate:"required"`
	Slot                           uint64                                   `json:"slot" validate:"required"`
	AdditionalWorkStaffInformation []AdditionalWorkStaffInformationResponse `json:"additionalWorkStaffInformation"`
}

type AdditionalWorkStaffInformationResponse struct {
	Id        uint64 `json:"id"`
	Date      string `json:"date"`
	Time      string `json:"time"`
	StaffName string `json:"staffName"`
	IsDone    bool   `json:"isDone"`
}

type AdditionalWorkDetailResponse struct {
	Id          uint64 `json:"id"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Time        string `json:"time"`
}

type AdditionalWorkStaffResponse struct {
	Id             uint64                       `json:"id"`
	IsDone         bool                         `json:"isDone"`
	AdditionalWork AdditionalWorkDetailResponse `json:"additionalWork"`
}

type DailyWorkStaffResponse struct {
	Id        uint64                  `json:"id"`
	IsDone    bool                    `json:"isDone"`
	DailyWork DailyWorkDetailResponse `json:"dailyWork"`
}

type DailyWorkListResponse struct {
	Role       RoleResponse `json:"role"`
	TotalWork  uint64       `json:"totalWork"`
	TotalStaff uint64       `json:"totalStaff"`
}

type AdditionalWorkListResponse struct {
	Id            uint64 `json:"id"`
	Date          string `json:"date"`
	Description   string `json:"description"`
	Location      string `json:"location"`
	RemainingSlot uint64 `json:"remainingSlot"`
	Status        string `json:"status"`
}

type GetAdditonalWorkFilter struct {
	Status string `query:"status"`
}

type UpdateAdditionalWorkStaffRequest struct {
	IsDone bool `json:"isDone"`
}

type UpdateDailyWorkStaffRequest struct {
	IsDone bool `json:"isDone"`
}
