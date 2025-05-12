package dto

import "github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/param"

type StaffResponse struct {
	Id           string       `json:"id"`
	Email        string       `json:"email"`
	Name         string       `json:"name"`
	PhotoProfile string       `json:"photoProfile"`
	PhoneNumber  string       `json:"phoneNumber"`
	Address      string       `json:"address"`
	Salary       string       `json:"salary"`
	Role         RoleResponse `json:"role"`
	CreatedAt    string       `json:"createdAt"`
}

type UpdateStaffRequest struct {
	Email        string `json:"email" validate:"required,email"`
	RoleId       uint64 `json:"roleId" validate:"required"`
	PhotoProfile string `json:"photoProfile" validate:"required"`
	Name         string `json:"name" validate:"required"`
	PhoneNumber  string `json:"phoneNumber" validate:"required"`
	Address      string `json:"address" validate:"required"`
	Salary       string `json:"salary" validate:"required"`
}

type GetStaffFilter struct {
	RoleId  string `query:"roleId"`
	Page    uint64 `query:"page"`
	Keyword string `query:"keyword"`
}

type StaffListResponse struct {
	Id           string       `json:"id"`
	Name         string       `json:"name"`
	PhotoProfile string       `json:"photoProfile"`
	Email        string       `json:"email"`
	Salary       string       `json:"salary"`
	PhoneNumber  string       `json:"phoneNumber"`
	Address      string       `json:"address"`
	Role         RoleResponse `json:"role"`
}

type StaffListPaginationResponse struct {
	TotalPage uint64              `json:"totalPage"`
	TotalData uint64              `json:"totalData"`
	Staffs    []StaffListResponse `json:"staffs"`
}

type StaffInformationResponse struct {
	TotalWorkHour float64 `json:"totalWorkHour"`
	KPIScore      float64 `json:"kpiScore"`
	TotalSalary   string  `json:"totalSalary"`
}

type KPIPerformanceResponse struct {
	Key   string  `json:"key"`
	Value float64 `json:"value"`
}

type StaffPresenceInformationResponse struct {
	TotalPresent    uint64 `json:"totalPresent"`
	TotalNotPresent uint64 `json:"totalNotPresent"`
}

type StaffWorkInformationResponse struct {
	TotalWorkDone    uint64 `json:"totalWorkDone"`
	TotalWorkNotDone uint64 `json:"totalWorkNotDone"`
}

type StaffSalaryInformationResponse struct {
	BaseSalary     string `json:"baseSalary"`
	OvertimeSalary string `json:"overtimeSalary"`
	BonusSalary    string `json:"bonusSalary"`
	Cashbon        string `json:"cashbon"`
	TotalSalary    string `json:"totalSalary"`
}

type GetStaffOverviewFilter struct {
	Month param.MonthParam `query:"month" validate:"required"`
	Year  uint64           `query:"year" validate:"required"`
}

type StaffOverviewResponse struct {
	StaffInformation         StaffInformationResponse         `json:"staffInformation"`
	KPIPerformances          []KPIPerformanceResponse         `json:"kpiPerformances"`
	StaffPresenceInformation StaffPresenceInformationResponse `json:"staffPresenceInformation"`
	StaffWorkInformation     StaffWorkInformationResponse     `json:"staffWorkInformation"`
	StaffSalaryInformation   StaffSalaryInformationResponse   `json:"staffSalaryInformation"`
}
