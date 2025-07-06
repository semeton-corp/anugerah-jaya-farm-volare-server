package mapper

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
)

func SupplierToResponse(supplier *entity.Supplier) dto.SupplierResponse {
	response := dto.SupplierResponse{
		Id:          supplier.Id,
		Name:        supplier.Name,
		PhoneNumber: supplier.PhoneNumber,
		Address:     supplier.Address,
	}

	items := make([]dto.ItemResponse, 0)
	for _, e := range supplier.SupplierItems {
		items = append(items, ItemToResponse(&e.Item))
	}

	response.Items = items

	return response
}

func SupplierToListResponse(supplier *entity.Supplier) dto.SupplierListResponse {
	response := dto.SupplierListResponse{
		Id:          supplier.Id,
		Name:        supplier.Name,
		PhoneNumber: supplier.PhoneNumber,
		Address:     supplier.Address,
	}

	return response
}
