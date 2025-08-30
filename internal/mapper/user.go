package mapper

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
)

func UserToResponse(user *entity.User) dto.UserResponse {
	response := dto.UserResponse{
		Id:           user.Id.String(),
		Name:         user.Name,
		Username:     user.Username,
		Email:        user.Email,
		Address:      user.Address,
		PhotoProfile: user.PhotoProfile,
		PhoneNumber:  user.PhoneNumber,
		Salary:       user.Salary.String(),
		CreatedAt:    user.CreatedAt.Format("02 Januari 2006"),
		Role:         RoleToResponse(&user.Role),
		Location:     LocationToResponse(&user.Location),
	}

	return response
}

func UserToListResponse(user *entity.User) dto.UserListResponse {
	return dto.UserListResponse{
		Id:           user.Id.String(),
		Name:         user.Name,
		Email:        user.Email,
		PhotoProfile: user.PhotoProfile,
		PhoneNumber:  user.PhoneNumber,
		Role: dto.RoleResponse{
			Id:   user.Role.Id,
			Name: user.Role.Name,
		},
	}
}

// Todo : KPI Status
func UserOverviewToListResponse(user *entity.User) dto.UserListOverviewResponse {
	return dto.UserListOverviewResponse{
		Id:                   user.Id.String(),
		Name:                 user.Name,
		Username:             user.Username,
		Email:                user.Email,
		PhotoProfile:         user.PhotoProfile,
		Role:                 RoleToResponse(&user.Role),
		SalaryRecommendation: user.Salary.String(),
		KpiStatus:            "",
	}
}
