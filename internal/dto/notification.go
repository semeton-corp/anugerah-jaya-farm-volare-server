package dto

type CreateNotificationRequest struct {
	UserId      *string `json:"userId"`
	StoreId     *uint64 `json:"storeId"`
	CageId      *uint64 `json:"cageId"`
	WarehouseId *uint64 `json:"warehouseId"`
	Description string  `json:"description" validate:"required"`
}

type NotificationResponse struct {
	Id          uint64 `json:"id"`
	Description string `json:"description"`
	IsMarked    bool   `json:"isMarked"`
}

type MarkNotificationRequest struct {
	Ids []uint64 `json:"ids"`
}

type GetNotificationFilter struct {
	UserId               string   `query:"userId"`
	StoreId              uint64   `query:"storeId"`
	CageId               uint64   `query:"cageId"`
	WarehouseId          uint64   `query:"warehouseId"`
	NotificationContexts []string `query:"notificationContexts"`
	IsMarked             *bool    `query:"isMarked"`
}
