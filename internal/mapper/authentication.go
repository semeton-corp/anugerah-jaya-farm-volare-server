package mapper

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
)

func AccountToResponse(data *entity.Account) dto.AccountResponse {
	return dto.AccountResponse{
		Id:           data.Id.String(),
		Email:        data.Email,
		PhotoProfile: data.PhotoProfile,
		Role:         RoleToResponse(&data.Role),
	}
}
