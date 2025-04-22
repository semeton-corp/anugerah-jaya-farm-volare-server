package dto

type StaffResponse struct {
	Id   string       `json:"id"`
	Name string       `json:"name"`
	Role RoleResponse `json:"role"`
}

type UpdateStaffRequest struct {
	Name string `json:"name" validate:"required"`
}
