package logger

import (
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DefaultTimeFamat = "06/01/02T15:04:05.000 MST"
)

func NewSugaredLogger(level, field string) *zap.SugaredLogger {
	//atom := zap.NewAtomicLevelAt(logLevel)
	config := NewDefaultConfig()
	config.Level = zap.NewAtomicLevelAt(nameToLevel(level))
	logger, _ := config.Build()
	return logger.Named(field).Sugar()
}

// NewLogger result not support Printf
func NewLogger(level, field string) *zap.Logger {
	//atom := zap.NewAtomicLevelAt(logLevel)
	config := NewDefaultConfig()
	config.Level = zap.NewAtomicLevelAt(nameToLevel(level))
	logger, _ := config.Build()
	return logger.Named(field)
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(DefaultTimeFamat))
}

func NewDefaultConfig() zap.Config {
	return zap.Config{
		Level:       zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Development: false,
		Encoding:    "console",
		EncoderConfig: zapcore.EncoderConfig{
			// Keys can be anything except the empty string.
			TimeKey:        "T",
			LevelKey:       "L",
			NameKey:        "N",
			CallerKey:      "C",
			MessageKey:     "M",
			StacktraceKey:  "S",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		//InitialFields:    map[string]interface{}{"serviceName": "wisdom_park"}, // 初始化字段，如：添加一个服务器名称
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

func nameToLevel(level string) zapcore.Level {
	switch level {
	case "DEBUG", "debug":
		return zapcore.DebugLevel
	case "INFO", "info":
		return zapcore.InfoLevel
	case "WARNING", "WARN", "warning", "warn":
		return zapcore.WarnLevel
	case "ERROR", "error":
		return zapcore.ErrorLevel
	case "DPANIC", "dpanic":
		return zapcore.DPanicLevel
	case "PANIC", "panic":
		return zapcore.PanicLevel
	case "FATAL", "fatal":
		return zapcore.FatalLevel

	case "NOTICE", "notice":
		return zapcore.InfoLevel // future
	case "CRITICAL", "critical":
		return zapcore.ErrorLevel // future

	default:
		fmt.Println("Unkown level using default level info")
		return zapcore.InfoLevel
	}
}

func FormatArgs(args ...interface{}) string { return strings.TrimSuffix(fmt.Sprintln(args...), "\n") }
