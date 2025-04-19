package mapper

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
)

func RoleToResponse(role *entity.Role) dto.RoleResponse {
	return dto.RoleResponse{
		Id:   role.Id,
		Name: role.Name,
	}
}
