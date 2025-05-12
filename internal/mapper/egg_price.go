package mapper

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
)

func EggPriceToResponse(eggPrice *entity.EggPrice) dto.EggPriceResponse {
	return dto.EggPriceResponse{
		Id:       eggPrice.Id,
		Category: eggPrice.Category,
		WarehouseItem: dto.WarehouseItemResponse{
			Id:       eggPrice.WarehouseItemId,
			Name:     eggPrice.WarehouseItem.Name,
			Unit:     eggPrice.WarehouseItem.Unit,
			Category: eggPrice.WarehouseItem.Category.String(),
		},
		Price: eggPrice.Price.String(),
	}
}

func EggPriceDiscountToResponse(eggPriceDiscount *entity.EggPriceDiscount) dto.EggPriceDiscountResponse {
	return dto.EggPriceDiscountResponse{
		Id:                     eggPriceDiscount.Id,
		Name:                   eggPriceDiscount.Name,
		MinimumTransactionUser: eggPriceDiscount.MinimumTransactionUser,
		TotalDiscount:          eggPriceDiscount.TotalDiscount,
	}
}
