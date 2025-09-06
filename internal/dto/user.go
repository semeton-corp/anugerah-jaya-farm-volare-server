package dto

import "github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/param"

type UserResponse struct {
	Id           string           `json:"id,omitempty"`
	Username     string           `json:"username,omitempty"`
	Email        string           `json:"email,omitempty"`
	Name         string           `json:"name,omitempty"`
	PhotoProfile string           `json:"photoProfile,omitempty"`
	PhoneNumber  string           `json:"phoneNumber,omitempty"`
	Address      string           `json:"address,omitempty"`
	Salary       string           `json:"salary,omitempty"`
	Role         RoleResponse     `json:"role,omitempty"`
	Location     LocationResponse `json:"location,omitzero"`
	CreatedAt    string           `json:"createdAt,omitempty"`
}

type UpdateUserRequest struct {
	Email        string `json:"email" validate:"required,email"`
	Username     string `json:"username" validate:"required"`
	LocationId   uint64 `json:"locationId" validate:"required"`
	RoleId       uint64 `json:"roleId" validate:"required"`
	PhotoProfile string `json:"photoProfile" validate:"required"`
	Name         string `json:"name" validate:"required"`
	PhoneNumber  string `json:"phoneNumber" validate:"required"`
	Address      string `json:"address" validate:"required"`
	Salary       string `json:"salary" validate:"required"`
}

type GetUserListFilter struct {
	RoleId         uint64   `query:"roleId"`
	LocationId     uint64   `query:"locationId"`
	ExcluseRoleIds []uint64 `query:"excludeRoleIds"`
}

type GetUserOverviewListFilter struct {
	Page           uint64 `query:"page"`
	Keyword        string `query:"keyword"`
	RoleId         uint64 `query:"roleId"`
	ExcludeRoleIds []uint64 `query:"excludeRoleIds"`
}

type UserListResponse struct {
	Id           string           `json:"id"`
	Name         string           `json:"name"`
	Email        string           `json:"email"`
	PhoneNumber  string           `json:"phoneNumber"`
	PhotoProfile string           `json:"photoProfile"`
	Role         RoleResponse     `json:"role"`
	Location     LocationResponse `json:"location"`
}

type UserListOverviewPaginationResponse struct {
	TotalPage uint64                     `json:"totalPage,omitempty"`
	TotalData uint64                     `json:"totalData,omitempty"`
	Users     []UserListOverviewResponse `json:"users"`
}

type UserInformationResponse struct {
	TotalWorkHour   float64 `json:"totalWorkHour"`
	WorkKpiScore    float64 `json:"workKpiScore"`
	ChickenKpiScore float64 `json:"chickenKpiScore"`
}

type KPIPerformanceResponse struct {
	Key             string  `json:"key"`
	WorkKpiScore    float64 `json:"workKpiScore"`
	ChickenKpiScore float64 `json:"chickenKpiScore"`
}

type UserPresenceInformationResponse struct {
	TotalPresent    uint64 `json:"totalPresent"`
	TotalNotPresent uint64 `json:"totalNotPresent"`
}

type UserWorkInformationResponse struct {
	TotalWorkDone    uint64 `json:"totalWorkDone"`
	TotalWorkNotDone uint64 `json:"totalWorkNotDone"`
}

type UserSalaryInformationResponse struct {
	BaseSalary           string `json:"baseSalary"`
	AdditionalWorkSalary string `json:"additionalWorkSalary"`
	BonusSalary          string `json:"bonusSalary"`
	CompentationSalary   string `json:"compentationSalary"`
	Cashbond             string `json:"cashbond"`
	IsPaid               bool   `json:"isPaid"`
	TotalSalary          string `json:"totalSalary"`
}

type GetUserOverviewFilter struct {
	Month param.MonthParam `query:"month" validate:"required"`
	Year  uint64           `query:"year" validate:"required"`
}

type UserOverviewResponse struct {
	User                    UserResponse                     `json:"user"`
	Placements              []string                         `json:"placements"`
	UseCashAdvances         []UserCashAdvanceSummaryResponse `json:"userCashAdvances"`
	UserInformation         UserInformationResponse          `json:"userInformation"`
	KPIPerformances         []KPIPerformanceResponse         `json:"kpiPerformances"`
	UserPresenceInformation UserPresenceInformationResponse  `json:"userPresenceInformation"`
	UserWorkInformation     UserWorkInformationResponse      `json:"userWorkInformation"`
	UserSalaryInformation   UserSalaryInformationResponse    `json:"userSalaryInformation"`
}

type UserListOverviewResponse struct {
	Id                   string       `json:"id"`
	Name                 string       `json:"name"`
	Username             string       `json:"username"`
	PhotoProfile         string       `json:"photoProfile"`
	Email                string       `json:"email"`
	SalaryRecommendation string       `json:"salaryRecommendation"`
	KpiStatus            string       `json:"kpiStatus"`
	Role                 RoleResponse `json:"role"`
}

type PerformanceGraphResponse struct {
	Key                   string  `json:"key"`
	KPIChickenPerformance float64 `json:"kpiChickenPerformance"`
	KPIUserPerformance    float64 `json:"kpiUserPerformance"`
}

type UserPerformanceSummaryResponse struct {
	TotalUser  uint64  `json:"totalUser"`
	KPIAll     float64 `json:"kpiAll"`
	KPIUser    float64 `json:"kpiUser"`
	KPIChicken float64 `json:"kpiChicken"`
}

type UserPerformanceOverviewResponse struct {
	UserPerformanceDetail UserPerformanceSummaryResponse `json:"userPerformanceSummary"`
	UserPerformanceGraphs []PerformanceGraphResponse     `json:"userPerformanceGraphs"`
}

type GetUserPerformanceOverviewFilter struct {
	LocationId uint64           `query:"locationId"`
	Month      param.MonthParam `query:"month" validate:"required"`
	Year       uint64           `query:"year" validate:"required"`
}

type UserSalaryGraphResponse struct {
	Key    string `json:"key"`
	Salary string `json:"kpiChickenPerformance"`
}
