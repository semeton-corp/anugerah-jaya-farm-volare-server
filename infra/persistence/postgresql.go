package persistence

import (
	"fmt"
	"time"

	_logger "github.com/semeton-corp/anugerah-jaya-farm-volare/infra/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New(log *zap.Logger) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		viper.GetString("database.host"),
		viper.GetString("database.user"),
		viper.GetString("database.pass"),
		viper.GetString("database.name"),
		viper.GetInt("database.port"),
		viper.GetString("database.sslmode"),
		viper.GetString("database.timezone"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:      _logger.NewZapGormLogger(log, logger.Warn),
		PrepareStmt: true,
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
		TranslateError: true,
	})

	if err != nil {
		zap.L().Panic("failed to connect database", zap.Error(err))
	}

	return db
}
