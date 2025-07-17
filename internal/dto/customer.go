package dto

type CustomerResponse struct {
	Id               uint64 `json:"id"`
	Name             string `json:"name"`
	PhoneNumber      string `json:"phoneNumber"`
	TotalTransaction uint64 `json:"totalTransaction,omitempty"`
}

type CreateCustomerRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
}
