package mapper

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
)

func StaffToResponse(staff *entity.Staff) dto.StaffResponse {
	return dto.StaffResponse{
		Id:   staff.Id.String(),
		Name: staff.Name,
		Role: dto.RoleResponse{
			Id:   staff.Account.Role.Id,
			Name: staff.Account.Role.Name,
		},
	}
}
