package mapper

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
)

func StaffToResponse(staff *entity.Staff) dto.StaffResponse {
	return dto.StaffResponse{
		Id:           staff.Id.String(),
		Name:         staff.Name,
		Email:        staff.Account.Email,
		Address:      staff.Address,
		PhotoProfile: staff.Account.PhotoProfile,
		PhoneNumber:  staff.PhoneNumber,
		Salary:       staff.Salary.String(), // just base salary
		CreatedAt:    staff.CreatedAt.Format("02 Januari 2006"),
		Role: dto.RoleResponse{
			Id:   staff.Account.Role.Id,
			Name: staff.Account.Role.Name,
		},
	}
}

// Note : for salary without adding salary for the lembur, bonus, kasbon
func StaffToListResponse(staff *entity.Staff) dto.StaffListResponse {
	return dto.StaffListResponse{
		Id:           staff.Id.String(),
		Name:         staff.Name,
		Email:        staff.Account.Email,
		Address:      staff.Address,
		PhotoProfile: staff.Account.PhotoProfile,
		PhoneNumber:  staff.PhoneNumber,
		Salary:       staff.Salary.String(), // just base salary
		Role: dto.RoleResponse{
			Id:   staff.Account.Role.Id,
			Name: staff.Account.Role.Name,
		},
	}
}
