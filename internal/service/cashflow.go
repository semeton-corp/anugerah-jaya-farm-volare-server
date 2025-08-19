package service

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"go.uber.org/zap"
)

type CashflowService struct {
	log              *zap.Logger
	storeService     IStoreService
	warehouseService IWarehouseService
	chickenService   IChickenService
	userService      IUserService
}

type ICashflowService interface {
}

func NewCashflowService(log *zap.Logger, storeService IStoreService, warehouseService IWarehouseService, chickenService IChickenService, userService IUserService) ICashflowService {
	return &CashflowService{
		log:              log,
		storeService:     storeService,
		warehouseService: warehouseService,
		chickenService:   chickenService,
		userService:      userService,
	}
}

func (s *CashflowService) GetIncomeOverview(filter dto.GetIncomeOverviewFilter) (dto.IncomeOverviewResponse, error) {
	return dto.IncomeOverviewResponse{}, nil
}
