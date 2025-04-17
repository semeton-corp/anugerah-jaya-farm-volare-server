package persistence

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(
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
	)
}
