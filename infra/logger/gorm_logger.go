package logger

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"go.uber.org/zap"
	gormlogger "gorm.io/gorm/logger"
)

type ZapGormLogger struct {
	ZapLogger     *zap.Logger
	LogLevel      gormlogger.LogLevel
	SlowThreshold time.Duration
}

type Config struct {
	SlowThreshold time.Duration
	LogLevel      gormlogger.LogLevel
}

func NewZapGormLogger(z *zap.Logger, config Config) gormlogger.Interface {
	if config.SlowThreshold == 0 {
		config.SlowThreshold = 200 * time.Millisecond // Default 200ms
	}

	return &ZapGormLogger{
		ZapLogger:     z,
		LogLevel:      config.LogLevel,
		SlowThreshold: config.SlowThreshold,
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
	if l.LogLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := []zap.Field{
		zap.Duration("elapsed", elapsed),
		zap.String("sql", formatSQL(sql)),
		zap.Int64("rows", rows),
	}

	queryType := extractQueryType(sql)
	if queryType != "" {
		fields = append(fields, zap.String("type", queryType))
	}

	tableName := extractTableName(sql)
	if tableName != "" {
		fields = append(fields, zap.String("table", tableName))
	}

	if err != nil {
		if l.LogLevel >= gormlogger.Error {
			fields = append(fields, zap.Error(err))
			l.ZapLogger.Error("Database query error", fields...)
		}
		return
	}

	if elapsed > l.SlowThreshold {
		if l.LogLevel >= gormlogger.Warn {
			fields = append(fields, zap.Duration("threshold", l.SlowThreshold))
			l.ZapLogger.Warn("Slow SQL query detected", fields...)
		}
		return
	}

	if l.LogLevel >= gormlogger.Info {
		l.ZapLogger.Debug("Database query executed", fields...)
	}
}

func convertToZapFields(data []interface{}) []zap.Field {
	var fields []zap.Field
	for i, item := range data {
		fields = append(fields, zap.Any(fmt.Sprintf("arg_%d", i), item))
	}
	return fields
}

func formatSQL(sql string) string {
	sql = strings.TrimSpace(sql)
	sql = regexp.MustCompile(`\s+`).ReplaceAllString(sql, " ")

	const maxLength = 1000
	if len(sql) > maxLength {
		return sql[:maxLength] + "... [truncated]"
	}

	return sql
}

func extractQueryType(sql string) string {
	sql = strings.TrimSpace(strings.ToUpper(sql))

	operations := []string{"SELECT", "INSERT", "UPDATE", "DELETE", "CREATE", "DROP", "ALTER", "TRUNCATE", "BEGIN", "COMMIT", "ROLLBACK"}

	for _, op := range operations {
		if strings.HasPrefix(sql, op) {
			return op
		}
	}

	return ""
}

func extractTableName(sql string) string {
	sql = strings.TrimSpace(sql)
	sqlUpper := strings.ToUpper(sql)

	if strings.HasPrefix(sqlUpper, "SELECT") {
		fromIdx := strings.Index(sqlUpper, " FROM ")
		if fromIdx != -1 {
			afterFrom := sql[fromIdx+6:] // +6 for " FROM "
			return extractFirstWord(afterFrom)
		}
	}

	if strings.HasPrefix(sqlUpper, "INSERT") {
		intoIdx := strings.Index(sqlUpper, " INTO ")
		if intoIdx != -1 {
			afterInto := sql[intoIdx+6:] // +6 for " INTO "
			return extractFirstWord(afterInto)
		}
	}

	if strings.HasPrefix(sqlUpper, "UPDATE") {
		parts := strings.Fields(sql)
		if len(parts) >= 2 {
			return extractFirstWord(parts[1])
		}
	}

	if strings.HasPrefix(sqlUpper, "DELETE") {
		fromIdx := strings.Index(sqlUpper, " FROM ")
		if fromIdx != -1 {
			afterFrom := sql[fromIdx+6:]
			return extractFirstWord(afterFrom)
		}
	}

	return ""
}

func extractFirstWord(s string) string {
	s = strings.TrimSpace(s)

	if strings.HasPrefix(s, "`") || strings.HasPrefix(s, "\"") || strings.HasPrefix(s, "'") {
		quote := s[0]
		endIdx := strings.IndexByte(s[1:], quote)
		if endIdx != -1 {
			return s[:endIdx+2]
		}
	}

	for i, r := range s {
		if r == ' ' || r == '(' || r == '\t' || r == '\n' {
			return s[:i]
		}
	}

	return s
}
