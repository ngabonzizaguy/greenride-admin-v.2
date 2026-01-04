package log

import (
	"fmt"
	"greenride/internal/config"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	// Service-specific loggers
	serviceLoggers = make(map[string]*logrus.Logger)
	loggerMutex    sync.RWMutex
)

const (
	DefaultLogger = ""
)

func Get() *logrus.Logger {
	return GetServiceLogger(DefaultLogger)
}

func GetServiceLogger(service string) *logrus.Logger {
	today := time.Now().Format(time.DateOnly)
	loggerKey := service
	if service != "" {
		loggerKey = fmt.Sprintf("%s-", loggerKey)
	}
	loggerKey = fmt.Sprintf("%s%s", loggerKey, today)
	return GetServiceLoggerWithoutDate(loggerKey)
}

func GetServiceLoggerWithoutDate(service string) *logrus.Logger {
	loggerKey := service
	loggerMutex.RLock()
	defer loggerMutex.RUnlock()
	if logger, exists := serviceLoggers[loggerKey]; exists {
		return logger
	}
	cfg := config.Get().Log
	logger := initLogger(cfg.Path, cfg.Level, service)
	serviceLoggers[loggerKey] = logger
	return logger
}

func initLogger(path string, level string, service string) (logger *logrus.Logger) {
	logger = logrus.New()
	// 如果没有配置日志路径，使用默认路径
	if path == "" {
		// 使用相对路径，确保在当前工作目录下创建日志
		path = "logs"
	}

	// 确保日志路径是绝对路径
	if !filepath.IsAbs(path) {
		currentDir, err := os.Getwd()
		if err == nil {
			path = filepath.Join(currentDir, path)
		}
	}

	// 如果日志路径不存在，则创建
	if err := os.MkdirAll(path, 0755); err != nil {
		log.Printf("Failed to create log directory: %v", err)
		// 降级到当前目录
		path = "logs"
		os.MkdirAll(path, 0755)
	}

	// 日志文件名格式：service.log
	logFileName := filepath.Join(path, fmt.Sprintf("%s.log", service))

	// 打开或创建日志文件
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		logger.SetOutput(os.Stdout)
	} else {
		// 设置多重输出，同时写入文件和标准输出
		mw := io.MultiWriter(os.Stdout, logFile)
		logger.SetOutput(mw)
	}

	// 设置日志格式为JSON
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	// 设置日志级别
	if level == "" {
		level = "info"
	}
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		logger.SetLevel(logrus.InfoLevel)
	} else {
		logger.SetLevel(lvl)
	}

	return
}

// Info logs info level message
func Info(args ...interface{}) {
	Get().Info(args...)
}

// Infof logs formatted info level message
func Infof(format string, args ...interface{}) {
	Get().Infof(format, args...)
}

// Error logs error level message
func Error(args ...interface{}) {
	Get().Error(args...)
}

// Errorf logs formatted error level message
func Errorf(format string, args ...interface{}) {
	Get().Errorf(format, args...)
}

// Warn logs warning level message
func Warn(args ...interface{}) {
	Get().Warn(args...)
}

// Warnf logs formatted warning level message
func Warnf(format string, args ...interface{}) {
	Get().Warnf(format, args...)
}

// Fatal logs fatal level message
func Fatal(args ...interface{}) {
	Get().Fatal(args...)
}

// Fatalf logs formatted fatal level message
func Fatalf(format string, args ...interface{}) {
	Get().Fatalf(format, args...)
}
