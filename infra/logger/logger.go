package logger

import (
	"time"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/constant"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Todo : fix log wiht not json type
func New() *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	var (
		log *zap.Logger
		err error
	)

	if viper.GetString("log.environment") == constant.LogEnvDevelopment {
		logCfg := zap.Config{
			Level:             zap.NewAtomicLevelAt(zap.DebugLevel),
			Encoding:          "json",
			EncoderConfig:     encoderCfg,
			OutputPaths:       []string{"stdout"},
			DisableStacktrace: true,
		}
		log, err = logCfg.Build()
	} else if viper.GetString("log.environment") == constant.LogEnvProduction {
		logCfg := zap.Config{
			Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
			Encoding:          "json",
			EncoderConfig:     encoderCfg,
			OutputPaths:       []string{"stdout"},
			DisableStacktrace: true,
		}
		log, err = logCfg.Build()
	} else {
		zap.L().Panic("invalid log environment")
	}

	if err != nil {
		zap.L().Panic("failed to create zap logger", zap.Error(err))
	}

	zap.ReplaceGlobals(log)
	return log
}
