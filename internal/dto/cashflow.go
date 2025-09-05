package dto

import (
	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/param"
)

type IncomePieResponse struct {
	WarehouseEggSalePercentage float64 `json:"warehouseEggSalePercentage"`
	StoreEggSalePercentage     float64 `json:"storeEggSalePercentage"`
	AfkirChickenSalePercentage float64 `json:"afkirChickenSalePercentage"`
	UserCashAdvancePercentage  float64 `json:"userCashAdvancePercentage"`
}

type IncomeListResponse struct {
	ParentId     uint64  `json:"parentId"`
	Id           uint64  `json:"id"`
	Date         string  `json:"date"`
	PlaceName    string  `json:"placeName"`
	Category     string  `json:"category"`
	ItemName     string  `json:"itemName"`
	ItemUnit     string  `json:"itemUnit"`
	Quantity     float64 `json:"quantity"`
	CustomerName string  `json:"customerName" `
	Nominal      string  `json:"nominal"`
	PaymentProof string  `json:"paymentProof"`
}

type IncomeResponse struct {
	ParentId            uint64  `json:"parentId"`
	Id                  uint64  `json:"id"`
	Date                string  `json:"date"`
	Time                string  `json:"time"`
	Category            string  `json:"category"`
	PlaceName           string  `json:"placeName"`
	CustomerName        string  `json:"customerName" `
	CustomerPhoneNumber string  `json:"customerPhoneNumber"`
	ItemName            string  `json:"itemName"`
	ItemUnit            string  `json:"itemUnit"`
	Quantity            float64 `json:"quantity"`
	Nominal             string  `json:"nominal"`
	PaymentType         string  `json:"paymentType"`
	TotalPrice          string  `json:"totalPrice"`
	PaymentMethod       string  `json:"paymentMethod"`
	InputBy             string  `json:"inputBy"`
	PaymentProof        string  `json:"paymentProof"`
}

type IncomeOverviewResponse struct {
	IncomePie IncomePieResponse    `json:"incomePie"`
	Incomes   []IncomeListResponse `json:"incomes"`
}

type GetIncomeOverviewFilter struct {
	Month          param.MonthParam `query:"month" validate:"required"`
	Year           uint64           `query:"year" validate:"required"`
	IncomeCategory string           `query:"category" validate:"required,incomeCategory"`
}

type ExpensePieResponse struct {
	StaffPercentage                        float64 `json:"staffPercentage"`
	OperationalPercentage                  float64 `json:"operationalPercentage"`
	ChikckenProcuremtnPercentage           float64 `json:"chickenProcurementPercentage"`
	WarehouseItemProcurementPercentage     float64 `json:"warehouseItemProcurementPercentage"`
	WarehouseItemCornProcurementPercentage float64 `json:"warehouseItemCornProcurementPercentage"`
	OtherPercentage                        float64 `json:"otherPercentage"`
	UserCashAdvancePercentage              float64 `json:"userCashAdvancePercentage"`
	TaxPercentage                          float64 `jsonn:"taxPercentage"`
}

type ExpenseListResponse struct {
	Id           uint64 `json:"id"`
	Date         string `json:"date"`
	Category     string `json:"category"`
	Name         string `json:"name"`
	PlaceName    string `json:"location"`
	Nominal      string `json:"nominal"`
	ReceiverName string `json:"receiverName"`
	PaymentProof string `json:"paymentProof"`
}

type ExpenseResponse struct {
	Id                  uint64 `json:"id"`
	Date                string `json:"date"`
	Time                string `json:"time"`
	Category            string `json:"category"`
	PlaceName           string `json:"placeName"`
	Name                string `json:"name"`
	ReceiverName        string `json:"receiverName"`
	ReceiverPhoneNumber string `json:"receiverPhoneNumber"`
	Nominal             string `json:"nominal"`
	PaymentMethod       string `json:"paymentMethod"`
	PaymentProof        string `json:"paymentProof"`
	InputBy             string `json:"inputBy"`
}

type ExpenseOverviewResponse struct {
	ExpensePie ExpensePieResponse    `json:"expensePie"`
	Expenses   []ExpenseListResponse `json:"expenses"`
}

type GetExpenseOverviewFilter struct {
	Month           param.MonthParam `query:"month" validate:"required"`
	Year            uint64           `query:"year" validate:"required"`
	ExpenseCategory string           `query:"category" validate:"required,expenseCategory"`
}

type GetReceivablesOverviewFilter struct {
	Month               param.MonthParam `query:"month" validate:"required"`
	Year                uint64           `query:"year" validate:"required"`
	ReceivablesCategory string           `query:"category" validate:"required,receivablesCategory"`
}

type GetCashflowSaleReportFilter struct {
	Year  uint64           `query:"year" validate:"required"`
	Month param.MonthParam `query:"month" validate:"required"`
}

type CreateExpenseRequest struct {
	ExpenseCategory     string `json:"expenseCategory" validate:"required,expenseCategory"`
	LocationId          uint64 `json:"locationId" validate:"required,min=1"`
	LocationType        string `json:"locationType" validate:"required,locationType"`
	PlaceId             uint64 `json:"placeId" validate:"required"`
	Name                string `json:"name" validate:"required"`
	ReceiverName        string `json:"receiverName" validate:"required"`
	ReceiverPhoneNumber string `json:"receiverPhoneNumber"`
	Nominal             string `json:"nominal" validate:"required"`
	PaymentMethod       string `json:"paymentMethod" validate:"required,paymentMethod"`
	PaymentProof        string `json:"PaymentProof" validate:"required"`
	Description         string `json:"description"`
}

type GetExpenseFilter struct {
	StartDate  param.DateParam `query:"startDate"`
	EndDate    param.DateParam `query:"endDate"`
	LocationId uint64          `query:"locationId"`
}

type CreateUserCashAdvanceRequest struct {
	UserId              string `json:"userId" validate:"required"`
	Nominal             string `json:"nominal" validate:"required"`
	DeadlinePaymentDate string `json:"deadlinePaymentDate" validate:"required"`
}

type UserCashAdvanceSummaryResponse struct {
	Id                            uint64 `json:"id"`
	DeadlinePaymentDate           string `json:"deadlinePaymentDate"`
	IsMoreThanDeadlinePaymentDate bool   `json:"isMoreThanDeadlinePaymentDate"`
	Nominal                       string `json:"nominal"`
	RemainingPayment              string `json:"remainingPayment"`
}

type UserCashAdvancePaymentResponse struct {
	Id            uint64 `json:"id"`
	Date          string `json:"date"`
	Nominal       string `json:"nominal"`
	Remaining     string `json:"remaining"`
	PaymentMethod string `json:"paymentMethod"`
	PaymentProof  string `json:"paymentProof"`
}

type UserCashAdvanceResponse struct {
	Id                      uint64                           `json:"id"`
	User                    UserListResponse                 `json:"user"`
	Nominal                 string                           `json:"nominal"`
	DeadlinePaymentDate     string                           `json:"deadlinePaymentDate"`
	PaymentStatus           string                           `json:"paymentStatus"`
	UserCashAdvancePayments []UserCashAdvancePaymentResponse `json:"payments"`
	RemainingPayment        string                           `json:"remainingPayment"`
}
type CreateUserCashAdvancePaymentRequest struct {
	UserCashAdvanceId uint64 `json:"userCashAdvanceId"`
	PaymentDate       string `json:"paymentDate" validate:"required"`
	Nominal           string `json:"nominal" validate:"required,number"`
	PaymentProof      string `json:"paymentProof" validate:"required,url"`
	PaymentMethod     string `json:"paymentMethod" validate:"required,paymentMethod"`
}

type GetUserCashAdvanceFilter struct {
	DeadlinePaymentStartDate param.DateParam            `query:"deadlinePaymentStartDate"`
	DeadlinePaymentEndDate   param.DateParam            `query:"deadlinePaymentEndDate"`
	UserId                   uuid.UUID                  `query:"userId"`
	PaymentStatus            param.PaymentStatusParam   `query:"paymentStatus"`
	PaymentStatuses          []param.PaymentStatusParam `query:"paymentStatuses"`
	LocationId               uint64
	StartDate                param.DateParam `query:"startDate"`
	EndDate                  param.DateParam `query:"endDate"`
}

type ReceivablesResponse struct {
	Id                    uint64                        `json:"id"`
	Date                  string                        `json:"date"`
	Time                  string                        `json:"time"`
	Category              string                        `json:"category"`
	PlaceName             string                        `json:"placeName"`
	Name                  string                        `json:"name"`
	PhoneNumber           string                        `json:"phoneNumber"`
	Nominal               string                        `json:"nominal"`
	RemainingPayment      string                        `json:"remainingPayment"`
	PaymentType           string                        `json:"paymentType"`
	PaymentStatus         string                        `json:"paymentStatus"`
	DeadlinePaymentDate   string                        `json:"deadlinePaymentDate"`
	InputBy               string                        `json:"inputBy"`
	ReceieveablesPayments []ReceievablesPaymentResponse `json:"payments"`
}

type ReceievablesPaymentResponse struct {
	Id            uint64 `json:"id"`
	Date          string `json:"date"`
	Nominal       string `json:"nominal"`
	Remaining     string `json:"remaining"`
	PaymentMethod string `json:"paymentMethod"`
	PaymentProof  string `json:"paymentProof"`
}

type ReceivablesListResponse struct {
	Id                  uint64 `json:"id"`
	DeadlinePaymentDate string `json:"deadlinePaymentDate"`
	Category            string `json:"category"`
	PlaceName           string `json:"placeName"`
	Name                string `json:"name"`
	PhoneNumber         string `json:"phoneNumber"`
	TotalNominal        string `json:"totalNominal"`
	RemainingPayment    string `json:"remainingPayment"`
	PaymentStatus       string `json:"paymentStatus"`
}

type ReceivablesPieResponse struct {
	UnpaidPercentage float64 `json:"unpaidPercentage"`
	PaidPercentage   float64 `json:"paidPercentage"`
}

type ReceivablesOverviewResponse struct {
	ReceivablesPie ReceivablesPieResponse    `json:"receivablesPie"`
	Receivables    []ReceivablesListResponse `json:"receivables"`
}

type UserSalaryPaymentResponse struct {
	Id                   uint64           `json:"id"`
	User                 UserListResponse `json:"user"`
	BaseSalary           string           `json:"baseSalary"`
	BonusSalary          string           `json:"bonusSalary"`
	CompentationSalary   string           `json:"compentationSalary"`
	AdditionalWorkSalary string           `json:"additionalWorkSalary"`
	PaymentProof         string           `json:"paymentProof"`
	PaymentMethod        string           `json:"paymentMethod"`
	IsPaid               bool             `json:"isPaid"`
}

type PayUserSalaryPaymentRequest struct {
	UserId                  string                                `json:"userId" validate:"required"`
	BaseSalary              string                                `json:"baseSalary" validate:"required"`
	BonusSalary             string                                `json:"bonusSalary" validate:"required"`
	CompentationSalary      string                                `json:"compentationSalary" validate:"required"`
	AdditionalWorkSalary    string                                `json:"additionalWorkSalary" validate:"required"`
	PaymentProof            string                                `json:"paymentProof" validate:"required,url"`
	PaymentMethod           string                                `json:"paymentMethod" validate:"required,paymentMethod"`
	UserCashAdvancePayments []CreateUserCashAdvancePaymentRequest `json:"userCashAdvancePayments"`
}

type DebtListResponse struct {
	Id                  uint64 `json:"id"`
	DeadlinePaymentDate string `json:"deadlinePaymentDate" validate:"required"`
	Category            string `json:"category"`
	PlaceName           string `json:"placeName"`
	TransactionName     string `json:"transactionName"`
	Name                string `json:"name"`
	PhoneNumber         string `json:"phoneNumber"`
	Nominal             string `json:"nominal"`
	RemainingPayment    string `json:"remainingPayment"`
	PaymentStatus       string `json:"paymentStatus"`
}

type DebtResponse struct {
	Id                  uint64                `jso:"id"`
	Date                string                `json:"date"`
	Time                string                `json:"time"`
	Category            string                `json:"category"`
	PlaceName           string                `json:"placeName"`
	TransactionName     string                `json:"transactionName"`
	Name                string                `json:"name"`
	PhoneNumber         string                `json:"phoneNumber"`
	DeadlinePaymentDate string                `json:"deadlinePaymentDate"`
	Nominal             string                `json:"nominal"`
	RemainingPayment    string                `json:"remainingPayment"`
	PaymentType         string                `json:"paymentType"`
	PaymentStatus       string                `json:"paymentStatus"`
	InputBy             string                `json:"inputBy"`
	DebtPayments        []DebtPaymentResponse `json:"payments"`
}

type DebtPaymentResponse struct {
	Id            uint64 `json:"id"`
	Date          string `json:"date"`
	Nominal       string `json:"nominal"`
	Remaining     string `json:"remaining"`
	PaymentMethod string `json:"paymentMethod"`
	PaymentProof  string `json:"paymentProof"`
}

type GetDebtOverviewFilter struct {
	Month        param.MonthParam `query:"month" validate:"required"`
	Year         uint64           `query:"year" validate:"required"`
	DebtCategory string           `query:"category" validate:"required,debtCategory"`
}

type DebtPieResponse struct {
	UnpaidPercentage float64 `json:"unpaidPercentage"`
	PaidPercentage   float64 `json:"paidPercentage"`
}

type DebtOverviewResponse struct {
	DebtPie DebtPieResponse    `json:"debtPie"`
	Debts   []DebtListResponse `json:"debts"`
}

type UserSalarySummaryResponse struct {
	TotalUser                uint64 `json:"totalUser"`
	TotalBaseSalary          string `json:"totalBaseSalary"`
	TotalAdditonalWorkSalary string `json:"totalAdditionalWorkSalary"`
	TotalBonusSalary         string `json:"totalBonusSalary"`
}

type GetUserSalarySummaryFilter struct {
	LocationId uint64           `query:"locationId"`
	Month      param.MonthParam `query:"month" validate:"required"`
	Year       uint64           `query:"year" validate:"required"`
}

type GetUserSalaryListFilter struct {
	LocationId uint64           `query:"locationId"`
	Month      param.MonthParam `query:"month" validate:"required"`
	Year       uint64           `query:"year" validate:"required"`
	Keyword    string           `query:"keyword"`
	RoleId     uint64           `query:"roleId"`
	Page       uint64           `query:"page"`
}

type UserSalaryListResponse struct {
	Id             uint64           `json:"id"`
	User           UserListResponse `json:"user"`
	SalaryInterval string           `json:"salaryInterval"`
	IsPaid         bool             `json:"isPaid"`
}

type UserSalaryListPaginationResponse struct {
	TotalData    uint64                   `json:"totalData,omitempty"`
	TotalPage    uint64                   `json:"totalPage,omitempty"`
	UserSalaries []UserSalaryListResponse `json:"userSalaries"`
}

type UserSalaryDetailResponse struct {
	Date                     string                           `json:"date,omitempty"`
	Time                     string                           `json:"time,omitempty"`
	User                     UserListResponse                 `json:"user"`
	SalaryMonth              string                           `json:"salaryMonth"`
	AdditionalWorkUsers      []AdditionalWorkUserResponse     `json:"additionalWorkUsers"`
	UserCashAdvanceSummaries []UserCashAdvanceSummaryResponse `json:"userCashAdvanceSummaries"`
	BaseSalary               string                           `json:"baseSalary"`
	BonusSalary              string                           `json:"bonusSalary"`
	CompentationSalary       string                           `json:"compentationSalary"`
	AdditionalWorkSalary     string                           `json:"additionalWorkSalary"`
}

type GetUserSalaryPaymentFilter struct {
	LocationId uint64          `query:"locationId"`
	StartDate  param.DateParam `query:"startDate"`
	EndDate    param.DateParam `query:"endDate"`
	IsPaid     *bool           `query:"isPaid"`
	Keyword    string          `query:"keyword"`
	RoleId     uint64          `query:"roleId"`
	Page       uint64          `query:"page"`
}

type CashflowSaleSummaryResponse struct {
	Income                string  `json:"income"`
	IsIncomeIncrease      bool    `json:"isIncomeIncrease"`
	IncomeDiffPercentage  float64 `json:"incomeDiffPercentage"`
	Profit                string  `json:"profit"`
	IsProfitIncrease      bool    `json:"isProfitIncrease"`
	ProfitDiffPercentage  float64 `json:"profitDiffPercentage"`
	Expense               string  `json:"expense"`
	IsExpenseIncrease     bool    `json:"isExpenseIncrease"`
	ExpenseDiffPercentage float64 `json:"expenseDiffPercentage"`
}

type CashflowSaleGraphResponse struct {
	Key     string `json:"key"`
	Income  string `json:"income"`
	Profit  string `json:"profit"`
	Expense string `json:"expense"`
}

type LocationSaleSummaryResponse struct {
	PlaceName     string `json:"placeName"`
	Income        string `json:"income"`
	Receieveables string `json:"receieveables"`
}

type LocationPieChartResponse struct {
	StorePercentage     float64 `json:"storePercentage"`
	WarehousePercentage float64 `json:"warehosuePercentage"`
}

type CashflowSaleOverviewResponse struct {
	CashflowSaleSummary CashflowSaleSummaryResponse   `json:"cashflowSaleSummary"`
	EggSaleSummary      EggSaleSummaryResponse        `json:"eggSaleSummary"`
	CashflowSaleGraphs  []CashflowSaleGraphResponse   `json:"cashflowSaleGraphs"`
	EggSaleGraphs       []EggSaleGraphResponse        `json:"eggSaleGraphs"`
	LocationSaleSummary []LocationSaleSummaryResponse `json:"locationSaleSummaries"`
	LocationPieChart    LocationPieChartResponse      `json:"locationPieChart"`
}

type GetCashflowSaleOverviewFilter struct {
	LocationId uint64           `query:"locationId"`
	Month      param.MonthParam `query:"month" validate:"required"`
	Year       uint64           `query:"year" validate:"required"`
	ItemId     uint64           `query:"itemId" validate:"required"`
}

type CashflowSummaryResponse struct {
	Income                    string  `json:"income"`
	IsIncomeIncrease          bool    `json:"isIncomeIncrease"`
	IncomeDiffPercentage      float64 `json:"incomeDiffPercentage"`
	Profit                    string  `json:"profit"`
	IsProfitIncrease          bool    `json:"isProfitIncrease"`
	ProfitDiffPercentage      float64 `json:"profitDiffPercentage"`
	Expense                   string  `json:"expense"`
	IsExpenseIncrease         bool    `json:"isExpenseIncrease"`
	ExpenseDiffPercentage     float64 `json:"expenseDiffPercentage"`
	Debt                      string  `json:"debt"`
	IsDebtIncrease            bool    `json:"isDebtIncrease"`
	DebtDiffPercentage        float64 `json:"debtDiffPercentage"`
	Cash                      string  `json:"cash"`
	IsCashIncrease            bool    `json:"isCashIncrease"`
	CashDiffPercentage        float64 `json:"cashDiffPercentage"`
	Receivables               string  `json:"receivables"`
	IsReceivablesIncrease     bool    `json:"isReceivablesIncrease"`
	ReceivablesDiffPercentage float64 `json:"receivablesDiffPercentage"`
}

type CashflowGraphResponse struct {
	Key     string `json:""`
	Income  string `json:"income"`
	Profit  string `json:"profit"`
	Expense string `json:"expense"`
	Cash    string `json:"cash"`
}

type EggSaleGraphResponse struct {
	Key   string  `json:"key"`
	Value float64 `json:"value"`
}

type EggSaleCashflowGraphResponse struct {
	Key              string `json:"key"`
	WarehouseEggSale string `json:"warehouseEggSale"`
	StoreEggSale     string `json:"storeEggSale"`
}

type GetCashflowOverviewFilter struct {
	Year       uint64 `query:"year" validate:"required"`
	LocationId uint64 `query:"locationId"`
}

type CashflowOverviewResponse struct {
	CashflowSummary       CashflowSummaryResponse        `json:"cashflowSummary"`
	CashflowGraphs        []CashflowGraphResponse        `json:"cashflowGraphs"`
	EggSaleCashflowGraphs []EggSaleCashflowGraphResponse `json:"eggSaleCashflowGraphs"`
}

type GetCashflowHistoryFilter struct {
	Year       uint64 `query:"year"`
	LocationId uint64 `query:"locationId"`
}

type GetUserCashAdvancePaymentFilter struct {
	StartDate  param.DateParam `query:"startDate"`
	EndDate    param.DateParam `query:"endDate"`
	LocationId uint64          `query:"locationId"`
}
