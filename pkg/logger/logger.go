package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// 日志级别
const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	FatalLevel = "fatal"
)

// 日志格式
const (
	TextFormat = "text"
	JSONFormat = "json"
)

// 内部日志级别
type logLevel int

const (
	debugLevel logLevel = iota
	infoLevel
	warnLevel
	errorLevel
	fatalLevel
)

var (
	// 默认日志格式为文本
	format = TextFormat
	// 默认日志级别为 info
	level = infoLevel
)

// 将字符串日志级别转换为内部日志级别
func parseLogLevel(lvl string) logLevel {
	switch strings.ToLower(lvl) {
	case DebugLevel:
		return debugLevel
	case InfoLevel:
		return infoLevel
	case WarnLevel:
		return warnLevel
	case ErrorLevel:
		return errorLevel
	case FatalLevel:
		return fatalLevel
	default:
		return infoLevel
	}
}

// InitTextLogger 初始化文本格式的日志记录器
func InitTextLogger(logLevel string) {
	format = TextFormat
	level = parseLogLevel(logLevel)
	log.SetFlags(log.Ldate | log.Ltime)
	log.SetOutput(os.Stdout)
}

// InitJSONLogger 初始化 JSON 格式的日志记录器
func InitJSONLogger(logLevel string) {
	format = JSONFormat
	level = parseLogLevel(logLevel)
	log.SetFlags(0) // 无标志，我们将手动格式化
	log.SetOutput(os.Stdout)
}

// 格式化日志消息
func formatLogMessage(lvl string, msg string) string {
	timestamp := time.Now().Format(time.RFC3339)
	
	if format == JSONFormat {
		return fmt.Sprintf(`{"timestamp":"%s","level":"%s","message":"%s"}`, 
			timestamp, lvl, escapeJSON(msg))
	}
	
	return fmt.Sprintf("[%s] %s: %s", timestamp, lvl, msg)
}

// 转义 JSON 字符串
func escapeJSON(s string) string {
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, `\`, `\\`)
	return s
}

// Debug 记录调试级别的日志
func Debug(msg string) {
	if level <= debugLevel {
		log.Println(formatLogMessage(DebugLevel, msg))
	}
}

// Info 记录信息级别的日志
func Info(msg string) {
	if level <= infoLevel {
		log.Println(formatLogMessage(InfoLevel, msg))
	}
}

// Warn 记录警告级别的日志
func Warn(msg string) {
	if level <= warnLevel {
		log.Println(formatLogMessage(WarnLevel, msg))
	}
}

// Error 记录错误级别的日志
func Error(msg string) {
	if level <= errorLevel {
		log.Println(formatLogMessage(ErrorLevel, msg))
	}
}

// Fatal 记录致命错误并终止程序
func Fatal(msg string) {
	if level <= fatalLevel {
		log.Println(formatLogMessage(FatalLevel, msg))
		os.Exit(1)
	}
} 