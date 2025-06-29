package mapper

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
)

func SupplierToResponse(supplier *entity.Supplier) dto.SupplierResponse {
	return dto.SupplierResponse{
		Id:            supplier.Id,
		WarehouseItem: ItemToResponse(&supplier.Item),
		Name:          supplier.Name,
		PhoneNumber:   supplier.PhoneNumber,
		Address:       supplier.Address,
	}
}
