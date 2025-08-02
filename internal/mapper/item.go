package mapper

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
)

func ItemToResponse(item *entity.Item) dto.ItemResponse {
	response := dto.ItemResponse{
		Id:       item.Id,
		Name:     item.Name,
		Category: item.Category.String(),
		Unit:     item.Unit,
	}

	if item.DailySpending.Valid {
		response.DailySpending = &item.DailySpending.Float64
	}

	return response
}

func ItemPriceToResponse(itemPrice *entity.ItemPrice) dto.ItemPriceResponse {
	return dto.ItemPriceResponse{
		Id:       itemPrice.Id,
		Category: itemPrice.Category,
		Item: dto.ItemResponse{
			Id:       itemPrice.ItemId,
			Name:     itemPrice.Item.Name,
			Unit:     itemPrice.Item.Unit,
			Category: itemPrice.Item.Category.String(),
		},
		Price: itemPrice.Price.String(),
	}
}

func ItemPriceDiscountToResponse(itemPriceDiscount *entity.ItemPriceDiscount) dto.ItemPriceDiscountResponse {
	return dto.ItemPriceDiscountResponse{
		Id:                     itemPriceDiscount.Id,
		Name:                   itemPriceDiscount.Name,
		MinimumTransactionUser: itemPriceDiscount.MinimumTransactionUser,
		TotalDiscount:          itemPriceDiscount.TotalDiscount,
	}
}
