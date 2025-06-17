package mapper

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
)

func LocationToResponse(location *entity.Location) dto.LocationResponse {
	return dto.LocationResponse{
		Id:   location.Id,
		Name: location.Name,
	}
}
