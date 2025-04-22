package persistence

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(
		&entity.Role{},
		&entity.Staff{},
		&entity.Account{},
		&entity.Location{},
		&entity.Cage{},
		&entity.ChickenMonitoring{},
		&entity.ChickenDiseaseMonitoring{},
		&entity.ChickenVaccineMonitoring{},
		&entity.EggMonitoring{},
		&entity.Warehouse{},
		&entity.WarehouseItem{},
		&entity.WarehouseStockItem{},
		&entity.Store{},
		&entity.StoreRequestItem{},
		&entity.StoreItem{},
		&entity.StoreSale{},
		&entity.StoreSalePayment{},
		&entity.DailyWork{},
		&entity.AdditionalWork{},
		&entity.AdditionalWorkStaff{},
		&entity.DailyWorkStaff{},
	)
}

func Rollback(db *gorm.DB) {
	db.Migrator().DropTable(
		&entity.Account{},
		&entity.Role{},
		&entity.Staff{},
		&entity.Location{},
		&entity.Cage{},
		&entity.ChickenMonitoring{},
		&entity.ChickenDiseaseMonitoring{},
		&entity.ChickenVaccineMonitoring{},
		&entity.EggMonitoring{},
		&entity.Warehouse{},
		&entity.WarehouseItem{},
		&entity.WarehouseStockItem{},
		&entity.Store{},
		&entity.StoreRequestItem{},
		&entity.StoreItem{},
		&entity.StoreSale{},
		&entity.StoreSalePayment{},
		&entity.DailyWork{},
		&entity.AdditionalWork{},
		&entity.AdditionalWorkStaff{},
		&entity.DailyWorkStaff{},
	)
}
