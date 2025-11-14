package logger

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/constant"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type DailyRotateWriter struct {
	logger      *lumberjack.Logger
	currentDate string
}

func (w *DailyRotateWriter) Write(p []byte) (n int, err error) {
	today := time.Now().Format("2006-01-02")

	if w.currentDate != "" && w.currentDate != today {
		w.logger.Rotate()
	}

	w.currentDate = today
	return w.logger.Write(p)
}

func (w *DailyRotateWriter) Sync() error {
	return nil
}

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

	useZipCompression := viper.GetBool("log.use-zip-compression")

	ll := &lumberjack.Logger{
		Filename:   "log/application.log",
		MaxSize:    100,                // Rotate when file reaches 100MB (or daily, whichever comes first)
		MaxBackups: 30,                 // Keep 30 backup files
		MaxAge:     30,                 // Keep logs for 30 days
		Compress:   !useZipCompression, // Use gzip if ZIP is disabled
		LocalTime:  true,               // Use local time for filenames
	}

	dailyWriter := &DailyRotateWriter{
		logger:      ll,
		currentDate: time.Now().Format("2006-01-02"),
	}

	// Start background task to compress old logs to ZIP format (if enabled)
	if useZipCompression {
		go compressOldLogsToZip()
	}

	logEnv := viper.GetString("log.environment")
	switch logEnv {
	case constant.LogEnvDevelopment:
		encoderCfg := zap.NewDevelopmentEncoderConfig()
		encoderCfg.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

		log = zap.New(zapcore.NewCore(zapcore.NewConsoleEncoder(encoderCfg), zap.CombineWriteSyncers(zapcore.Lock(os.Stdout), zapcore.AddSync(dailyWriter)), zapLevel), zap.AddStacktrace(zapcore.PanicLevel))
	case constant.LogEnvProduction:
		encoderCfg := zap.NewProductionEncoderConfig()
		encoderCfg.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

		log = zap.New(zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), zap.CombineWriteSyncers(zapcore.AddSync(dailyWriter)), zapLevel), zap.AddStacktrace(zapcore.PanicLevel))
	default:
		zap.L().Panic("invalid log environment")
	}

	zap.ReplaceGlobals(log)
	return log
}

func compressOldLogsToZip() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	compressLogFiles()

	for range ticker.C {
		compressLogFiles()
	}
}

func compressLogFiles() {
	logDir := "log"

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		return
	}

	files, err := os.ReadDir(logDir)
	if err != nil {
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filename := file.Name()

		// Only compress rotated log files (those with timestamp suffix)
		// Format: application-YYYY-MM-DD*.log (but not application.log)
		if strings.HasPrefix(filename, "application-") &&
			strings.HasSuffix(filename, ".log") &&
			!strings.HasSuffix(filename, ".zip") {

			filePath := filepath.Join(logDir, filename)

			info, err := file.Info()
			if err != nil {
				continue
			}

			if time.Since(info.ModTime()) > 1*time.Hour {
				if err := compressFileToZip(filePath); err == nil {
					os.Remove(filePath)
				}
			}
		}
	}
}

func compressFileToZip(filePath string) error {
	sourceFile, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	fileInfo, err := sourceFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	zipPath := filePath + ".zip"
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	header := &zip.FileHeader{
		Name:   filepath.Base(filePath),
		Method: zip.Deflate, // Use Deflate compression algorithm (not just Store)
	}
	header.Modified = fileInfo.ModTime()

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("failed to create file in zip: %w", err)
	}

	if _, err := io.Copy(writer, sourceFile); err != nil {
		return fmt.Errorf("failed to copy content to zip: %w", err)
	}

	return nil
}
