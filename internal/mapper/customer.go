package mapper

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
)

func CustomerToResponse(customer *entity.Customer) dto.CustomerResponse {
	return dto.CustomerResponse{
		Id:               customer.Id,
		Name:             customer.Name,
		PhoneNumber:      customer.PhoneNumber,
		TotalTransaction: uint64(len(customer.StoreSales)) + uint64(len(customer.WarehouseSales)),
	}
}
