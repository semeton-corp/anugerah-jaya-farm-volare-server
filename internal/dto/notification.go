package dto

type CreateNotificationRequest struct {
	UserId      *string `json:"userId"`
	StoreId     *uint64 `json:"storeId"`
	CageId      *uint64 `json:"cageId"`
	WarehouseId *uint64 `json:"warehouseId"`
	Description string  `json:"description" validate:"required"`
}

type NotificationResponse struct {
	Id                   uint64   `json:"id"`
	Description          string   `json:"description"`
	IsMarked             bool     `json:"isMarked"`
	NotificationContexts []string `json:"notificationContexts"`
}

type MarkNotificationRequest struct {
	Ids []uint64 `json:"ids"`
}

type GetNotificationFilter struct {
	UserIds              []string `query:"userIds"`
	StoreIds             []uint64 `query:"storeIds"`
	CageIds              []uint64 `query:"cageIds"`
	WarehouseIds         []uint64 `query:"warehouseIds"`
	NotificationContexts []string `query:"notificationContexts"`
	IsMarked             *bool    `query:"isMarked"`
}
