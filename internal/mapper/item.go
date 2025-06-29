package mapper

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
)

func ItemToResponse(warehouseItem *entity.Item) dto.ItemResponse {
	return dto.ItemResponse{
		Id:       warehouseItem.Id,
		Name:     warehouseItem.Name,
		Category: warehouseItem.Category.String(),
		Unit:     warehouseItem.Unit,
	}
}

func ItemPriceToResponse(eggPrice *entity.ItemPrice) dto.ItemPriceResponse {
	return dto.ItemPriceResponse{
		Id:       eggPrice.Id,
		Category: eggPrice.Category,
		Item: dto.ItemResponse{
			Id:       eggPrice.ItemId,
			Name:     eggPrice.Item.Name,
			Unit:     eggPrice.Item.Unit,
			Category: eggPrice.Item.Category.String(),
		},
		Price: eggPrice.Price.String(),
	}
}

func ItemPriceDiscountToResponse(eggPriceDiscount *entity.ItemPriceDiscount) dto.ItemPriceDiscountResponse {
	return dto.ItemPriceDiscountResponse{
		Id:                     eggPriceDiscount.Id,
		Name:                   eggPriceDiscount.Name,
		MinimumTransactionUser: eggPriceDiscount.MinimumTransactionUser,
		TotalDiscount:          eggPriceDiscount.TotalDiscount,
	}
}
