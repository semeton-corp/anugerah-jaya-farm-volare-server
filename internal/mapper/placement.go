package mapper

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
)

func CagePlacementToResponse(data *entity.CagePlacement) dto.CagePlacementResponse {
	return dto.CagePlacementResponse{
		User: UserToListResponse(&data.User),
		Cage: CageToResponse(&data.Cage),
	}
}

func StorePlacementToResponse(data *entity.StorePlacement) dto.StorePlacementResponse {
	return dto.StorePlacementResponse{
		User:  UserToListResponse(&data.User),
		Store: StoreToResponse(&data.Store),
	}
}

func WarehousePlacementToResponse(data *entity.WarehousePlacement) dto.WarehousePlacementResponse {
	return dto.WarehousePlacementResponse{
		User:      UserToListResponse(&data.User),
		Warehouse: WarehouseToResponse(&data.Warehouse),
	}
}
