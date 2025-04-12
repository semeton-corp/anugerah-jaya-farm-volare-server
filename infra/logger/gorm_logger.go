package logger

import (
	"context"
	"time"

	"go.uber.org/zap"
	gormlogger "gorm.io/gorm/logger"
)

type ZapGormLogger struct {
	ZapLogger *zap.Logger
	LogLevel  gormlogger.LogLevel
}

func NewZapGormLogger(z *zap.Logger, level gormlogger.LogLevel) gormlogger.Interface {
	return &ZapGormLogger{
		ZapLogger: z,
		LogLevel:  level,
	}
}

func (l *ZapGormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

func (l *ZapGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		l.ZapLogger.Info(msg, convertToZapFields(data)...)
	}
}

func (l *ZapGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		l.ZapLogger.Warn(msg, convertToZapFields(data)...)
	}
}

func (l *ZapGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		l.ZapLogger.Error(msg, convertToZapFields(data)...)
	}
}

func (l *ZapGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := []zap.Field{
		zap.Duration("elapsed", elapsed),
		zap.String("sql", sql),
		zap.Int64("rows", rows),
	}

	if err != nil && l.LogLevel >= gormlogger.Error {
		fields = append(fields, zap.Error(err))
		l.ZapLogger.Error("GORM SQL ERROR", fields...)
	} else if l.LogLevel >= gormlogger.Info {
		l.ZapLogger.Info("GORM SQL", fields...)
	}
}

func convertToZapFields(data []interface{}) []zap.Field {
	var fields []zap.Field
	for _, item := range data {
		fields = append(fields, zap.Any("data", item))
	}
	return fields
}
