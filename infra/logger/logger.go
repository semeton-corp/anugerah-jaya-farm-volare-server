package logger

import (
	"os"
	"time"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/constant"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func New() *zap.Logger {
	var (
		zapLevel zapcore.Level
		log      *zap.Logger
	)

	logLevel := viper.GetString("log.level")
	switch logLevel {
	case constant.LogLevelDebug:
		zapLevel = zap.DebugLevel
	case constant.LogLevelInfo:
		zapLevel = zap.InfoLevel
	case constant.LogLevelError:
		zapLevel = zap.ErrorLevel
	case constant.LogLevelWarn:
		zapLevel = zap.WarnLevel
	}

	ll := lumberjack.Logger{
		Filename:   "log/application.log",
		MaxSize:    1024,
		MaxBackups: 30,
		MaxAge:     30,
		Compress:   true,
	}

	logEnv := viper.GetString("log.environment")
	switch logEnv {
	case constant.LogEnvDevelopment:
		encoderCfg := zap.NewDevelopmentEncoderConfig()
		encoderCfg.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

		log = zap.New(zapcore.NewCore(zapcore.NewConsoleEncoder(encoderCfg), zap.CombineWriteSyncers(zapcore.Lock(os.Stdout), zapcore.AddSync(&ll)), zapLevel), zap.AddStacktrace(zapcore.PanicLevel))
	case constant.LogEnvProduction:
		encoderCfg := zap.NewProductionEncoderConfig()
		encoderCfg.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

		log = zap.New(zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), zap.CombineWriteSyncers(zapcore.AddSync(&ll)), zapLevel), zap.AddStacktrace(zapcore.PanicLevel))
	default:
		zap.L().Panic("invalid log environment")
	}

	zap.ReplaceGlobals(log)
	return log
}
