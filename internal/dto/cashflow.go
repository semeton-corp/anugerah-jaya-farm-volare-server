package dto

import (
	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/param"
)

type IncomePieResponse struct {
	WarehouseEggSalePercentage float64 `json:"warehouseEggSalePercentage"`
	StoreEggSalePercentage     float64 `json:"storeEggSalePercentage"`
	AfkirChickenSalePercentage float64 `json:"afkirChickenSalePercentage"`
}

type IncomeListResponse struct {
	ParentId     uint64 `json:"parentId"`
	Id           uint64 `json:"id"`
	Date         string `json:"date"`
	PlaceName    string `json:"placeName"`
	Category     string `json:"category"`
	ItemName     string `json:"itemName"`
	ItemUnit     string `json:"itemUnit"`
	Quantity     string `json:"quantity"`
	CustomerName string `json:"customerName" `
	Nominal      string `json:"nominal"`
	PaymentProof string `json:"paymentProof"`
}

type IncomeResponse struct {
	Id                  uint64 `json:"id"`
	Date                string `json:"date"`
	Time                string `json:"time"`
	Category            string `json:"category"`
	PlaceName           string `json:"placeName"`
	CustomerName        string `json:"customerName" `
	CustomerPhoneNumber string `json:"customerPhoneNumber"`
	ItemName            string `json:"itemName"`
	ItemUnit            string `json:"itemUnit"`
	Quantity            string `json:"quantity"`
	Nominal             string `json:"nominal"`
	PaymentType         string `json:"paymentType"`
	TotalPrice          string `json:"totalPrice"`
	PaymentMethod       string `json:"paymentMethod"`
	InputBy             string `json:"inputBy"`
	PaymentProof        string `json:"paymentProof"`
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
	StaffPercentage                float64 `json:"staffPercentage"`
	OperationalPercentage          float64 `json:"operationalPercentage"`
	ChikckenProcuremtnPercentage   float64 `json:"chickenProcurementPercentage"`
	WarehouseProcurementPercentage float64 `json:"warehouseItemProcurementPercentage"`
	OtherPercentage                float64 `json:"otherPercentage"`
}

type ExpenseListResponse struct {
	Id           uint64 `json:"id"`
	Date         string `json:"date"`
	Category     string `json:"category"`
	Name         string `json:"name"`
	PlaceName    string `json:"location"`
	Nominal      string `json:"string"`
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
	InputBy             string `json:"inputBy "`
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
	Month                 param.MonthParam `query:"month" validate:"required"`
	Year                  uint64           `query:"year" validate:"required"`
	ReceieveablesCategory string           `query:"category" validate:"required,receieveablesCategory"`
}

type GetSaleCashflowFilter struct {
	Year  uint64           `query:"year"`
	Month param.MonthParam `query:"month"`
}

type CreateExpenseRequest struct {
	ExpenseCategory     string `json:"expenseCategory" validate:"required,expenseCategory"`
	LocationId          uint64 `json:"locationId" validate:"required,min=1"`
	LocationType        string `json:"locationType" validate:"required"`
	PlaceId             uint64 `json:"placeId" validate:"required"`
	Name                string `json:"name" validate:"required"`
	ReceiverName        string `json:"receiverName" validate:"required"`
	ReceiverPhoneNumber string `json:"receiverPhoneNumber"`
	Nominal             string `json:"nominal" validate:"required"`
	PaymentMethod       string `json:"paymentMethod" validate:"required,paymentMethod"`
	PaymentProof        string `json:"paymentProof" validate:"required"`
	Description         string `json:"description"`
}

type GetExpenseFilter struct {
	StartDate param.DateParam `query:"startDate"`
	EndDate   param.DateParam `query:"endDate"`
}

type CreateCashAdvanceRequest struct {
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
	PaymentStatus           string                           `json:"deadlinePaymentStatus"`
	UserCashAdvancePayments []UserCashAdvancePaymentResponse `json:"payments"`
	RemainingPayment        string                           `json:"remainingPayment"`
}
type CreateUsereCashAdvancePaymentRequest struct {
	PaymentDate   string `json:"paymentDate" validate:"required"`
	Nominal       string `json:"nominal" validate:"required,number"`
	PaymentProof  string `json:"paymentProof" validate:"required,url"`
	PaymentMethod string `json:"paymentMethod" validate:"required,paymentMethod"`
}

type GetUserCashAdvanceFilter struct {
	DeadlinePaymentStartDate param.DateParam            `query:"deadlinePaymentStartDate"`
	DeadlinePaymentEndDate   param.DateParam            `query:"deadlinePaymentEndDate"`
	UserId                   uuid.UUID                  `query:"userId"`
	PaymentStatus            param.PaymentStatusParam   `query:"paymentStatus"`
	PaymentStatuses          []param.PaymentStatusParam `query:"paymentStatuses"`
}

type ReceiveablesResponse struct {
	Id                    uint64                         `json:"id"`
	Date                  string                         `json:"date"`
	Time                  string                         `json:"time"`
	Category              string                         `json:"category"`
	PlaceName             string                         `json:"placeName"`
	Name                  string                         `json:"name"`
	PhoneNumber           string                         `json:"phoneNumber"`
	RemainingPayment      string                         `json:"remainingPayment"`
	PaymentType           string                         `json:"paymentType"`
	PaymentStatus         string                         `json:"paymentStatus"`
	DeadlinePaymentDate   string                         `json:"deadlinePaymentDate"`
	InputBy               string                         `json:"inputBy"`
	ReceieveablesPayments []ReceieveablesPaymentResponse `json:"payments"`
}

type ReceieveablesPaymentResponse struct {
	Id            uint64 `json:"id"`
	Date          string `json:"date"`
	Nominal       string `json:"nominal"`
	Remaining     string `json:"remaining"`
	PaymentMethod string `json:"paymentMethod"`
	PaymentProof  string `json:"paymentProof"`
}

type ReceiveablesListResponse struct {
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

type ReceiveablesPieResponse struct {
	UnpaidPercentage  float64 `json:"unpaidPercentage"`
	PaidPercentage    float64 `json:"paidPercentage"`
	NotPaidPercentage float64 `json:"notPaidPercentage"`
}

type ReceievablesOverviewResponse struct {
	ReceivablesPie ReceiveablesPieResponse    `json:"receivablesPie"`
	Receivables    []ReceiveablesListResponse `json:"receivables"`
}

type UserSalaryPaymentResponse struct {
}

type PayUserSalaryPaymentRequest struct {
}

type DebtListResponse struct {
}

type DebtResponse struct {
}

type DebtPaymentResponse struct {
}

type GetDebtOverviewFilter struct {
	Month        param.MonthParam `query:"month" validate:"required"`
	Year         uint64           `query:"year" validate:"required"`
	DebtCategory string           `query:"category" validate:"required,debtCategory"`
}

type DebtPie struct{
	
}

type DebtOverviewResponse struct {
	Debts []DebtListResponse `json:"debts"`
}
