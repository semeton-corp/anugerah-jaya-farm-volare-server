package mapper

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
)

func UserToResponse(user *entity.User) dto.UserResponse {
	return dto.UserResponse{
		Id:           user.Id.String(),
		Name:         user.Name,
		Username:     user.Username,
		Email:        user.Email,
		Address:      user.Address,
		PhotoProfile: user.PhotoProfile,
		PhoneNumber:  user.PhoneNumber,
		Salary:       user.Salary.String(), //Note : just base salary
		CreatedAt:    user.CreatedAt.Format("02 Januari 2006"),
		Role:         RoleToResponse(&user.Role),
		Location:     LocationToResponse(&user.Location),
	}
}

// Note : for salary without adding salary for the lembur, bonus, kasbon
func UserToListResponse(user *entity.User) dto.UserListResponse {
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
