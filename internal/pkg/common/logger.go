package common

import (
	"fmt"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Logger 全局日誌實例
	Logger *zap.Logger

	// 定義日誌級別的顏色
	levelColors = map[zapcore.Level]string{
		zapcore.DebugLevel: "\033[36m", // 青色
		zapcore.InfoLevel:  "\033[32m", // 綠色
		zapcore.WarnLevel:  "\033[33m", // 黃色
		zapcore.ErrorLevel: "\033[31m", // 紅色
		zapcore.FatalLevel: "\033[35m", // 紫色
	}
	resetColor = "\033[0m"
)

// 自定義編碼器配置
func getEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "", // 移除 logger 名稱
		CallerKey:      "", // 移除調用者信息
		MessageKey:     "msg",
		StacktraceKey:  "", // 移除堆棧跟踪
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    customLevelEncoder,
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   nil, // 移除調用者編碼器
	}
}

// 自定義時間格式
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("15:04:05.000")) // 添加毫秒級別的時間戳
}

// 自定義級別編碼器（添加顏色）
func customLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	color := levelColors[l]
	level := l.String()
	// 統一級別顯示長度
	switch l {
	case zapcore.DebugLevel:
		level = "DBG"
	case zapcore.InfoLevel:
		level = "INF"
	case zapcore.WarnLevel:
		level = "WRN"
	case zapcore.ErrorLevel:
		level = "ERR"
	case zapcore.FatalLevel:
		level = "FAT"
	}
	enc.AppendString(color + level + resetColor)
}

// InitLogger 初始化日誌系統
func InitLogger(debug bool) error {
	// 設置日誌級別
	level := zapcore.InfoLevel
	if debug {
		level = zapcore.DebugLevel
	}

	// 創建日誌目錄
	if err := os.MkdirAll("logs", 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// 創建日誌文件
	logFile, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	// 創建多個輸出目標
	fileWriter := zapcore.AddSync(logFile)
	consoleWriter := zapcore.AddSync(os.Stdout)

	// 創建多個核心
	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(getEncoderConfig()),
		fileWriter,
		level,
	)
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(getEncoderConfig()),
		consoleWriter,
		level,
	)

	// 合併多個核心
	core := zapcore.NewTee(fileCore, consoleCore)

	// 創建 logger，移除一些默認字段
	Logger = zap.New(core,
		zap.AddCallerSkip(1),
		zap.Fields(
			zap.String("service", "recipe-generator"),
		),
	)

	// 替換全局 logger
	zap.ReplaceGlobals(Logger)

	return nil
}

// LogRequest 記錄 HTTP 請求
func LogRequest(method, path string, status int, duration time.Duration, requestID string) {
	if status >= 400 {
		Logger.Error("Request failed",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", status),
			zap.Duration("duration", duration),
		)
		return
	}
	Logger.Info("Request OK",
		zap.String("method", method),
		zap.String("path", path),
		zap.Duration("duration", duration),
	)
}

// LogInfo 記錄信息日誌
func LogInfo(msg string, fields ...zap.Field) {
	// 過濾掉包含圖片數據的字段
	filteredFields := make([]zap.Field, 0, len(fields))
	for _, field := range fields {
		if field.Key == "image" || strings.Contains(field.Key, "image_data") || strings.Contains(field.Key, "base64") {
			continue
		}
		filteredFields = append(filteredFields, field)
	}
	Logger.Info(msg, filteredFields...)
}

// LogError 記錄錯誤日誌
func LogError(msg string, fields ...zap.Field) {
	// 過濾掉包含圖片數據的字段
	filteredFields := make([]zap.Field, 0, len(fields))
	for _, field := range fields {
		if field.Key == "image" || strings.Contains(field.Key, "image_data") || strings.Contains(field.Key, "base64") {
			continue
		}
		filteredFields = append(filteredFields, field)
	}
	Logger.Error(msg, filteredFields...)
}

// LogWarn 記錄警告日誌
func LogWarn(msg string, fields ...zap.Field) {
	// 過濾掉包含圖片數據的字段
	filteredFields := make([]zap.Field, 0, len(fields))
	for _, field := range fields {
		if field.Key == "image" || strings.Contains(field.Key, "image_data") || strings.Contains(field.Key, "base64") {
			continue
		}
		filteredFields = append(filteredFields, field)
	}
	Logger.Warn(msg, filteredFields...)
}

// LogDebug 記錄調試日誌
func LogDebug(msg string, fields ...zap.Field) {
	// 過濾掉包含圖片數據的字段
	filteredFields := make([]zap.Field, 0, len(fields))
	for _, field := range fields {
		if field.Key == "image" || strings.Contains(field.Key, "image_data") || strings.Contains(field.Key, "base64") {
			continue
		}
		filteredFields = append(filteredFields, field)
	}
	Logger.Debug(msg, filteredFields...)
}

// LogFatal 記錄致命錯誤日誌
func LogFatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}

// Sync 同步日誌緩衝
func Sync() {
	if Logger != nil {
		_ = Logger.Sync()
	}
}

// LogCacheHit 記錄緩存命中
func LogCacheHit(cacheType, key string) {
	LogInfo("Cache Hit", zap.String("type", cacheType))
}

// LogCacheMiss 記錄緩存未命中
func LogCacheMiss(cacheType, key string) {
	LogInfo("Cache Miss", zap.String("type", cacheType))
}

// LogAICall 記錄 AI 調用
func LogAICall(prompt string, duration time.Duration, err error, requestID string) {
	if err != nil {
		LogError("AI Call Failed",
			zap.Error(err),
			zap.Duration("duration", duration),
		)
		return
	}
	LogInfo("AI Call OK",
		zap.Duration("duration", duration),
	)
}

// LogImageProcessing 記錄圖片處理相關的日誌
func LogImageProcessing(level string, msg string, fields ...zap.Field) {
	// 過濾掉包含圖片數據的字段
	filteredFields := make([]zap.Field, 0, len(fields))
	for _, field := range fields {
		if field.Key == "image" ||
			strings.Contains(field.Key, "image_data") ||
			strings.Contains(field.Key, "base64") ||
			field.Key == "has_image" {
			continue
		}
		filteredFields = append(filteredFields, field)
	}

	// 根據日誌級別記錄
	switch level {
	case "info":
		LogInfo(msg, filteredFields...)
	case "error":
		LogError(msg, filteredFields...)
	case "warn":
		LogWarn(msg, filteredFields...)
	default:
		LogInfo(msg, filteredFields...)
	}
}
