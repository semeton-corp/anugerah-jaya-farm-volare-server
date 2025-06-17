package mapper

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
)

func StaffToResponse(user *entity.User) dto.UserResponse {
	return dto.UserResponse{
		Id:           user.Id.String(),
		Name:         user.Name,
		Email:        user.Email,
		Address:      user.Address,
		PhotoProfile: user.PhotoProfile,
		PhoneNumber:  user.PhoneNumber,
		Salary:       user.Salary.String(), // just base salary
		CreatedAt:    user.CreatedAt.Format("02 Januari 2006"),
		Role: dto.RoleResponse{
			Id:   user.Role.Id,
			Name: user.Role.Name,
		},
	}
}

// Note : for salary without adding salary for the lembur, bonus, kasbon
func StaffToListResponse(user *entity.User) dto.UserListResponse {
	return dto.UserListResponse{
		Id:           user.Id.String(),
		Name:         user.Name,
		Email:        user.Email,
		Address:      user.Address,
		PhotoProfile: user.PhotoProfile,
		PhoneNumber:  user.PhoneNumber,
		Salary:       user.Salary.String(), // just base salary
		Role: dto.RoleResponse{
			Id:   user.Role.Id,
			Name: user.Role.Name,
		},
	}
}
