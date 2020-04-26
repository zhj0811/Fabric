package logging

import (
	"fmt"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestZap(t *testing.T) {
	atom := zap.NewAtomicLevel()

	atom.SetLevel(zap.DebugLevel)

	config := zap.Config{
		Level:       atom,
		Development: false,
		Encoding:    "console",
		EncoderConfig: zapcore.EncoderConfig{
			// Keys can be anything except the empty string.
			TimeKey:       "T",
			LevelKey:      "L",
			NameKey:       "N",
			CallerKey:     "C",
			MessageKey:    "M",
			StacktraceKey: "S",
			LineEnding:    zapcore.DefaultLineEnding,
			EncodeLevel:   zapcore.CapitalColorLevelEncoder,
			//EncodeLevel:   zapcore.CapitalLevelEncoder,
			//EncodeTime:    zapcore.ISO8601TimeEncoder,
			EncodeTime:     TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		//InitialFields:    map[string]interface{}{"serviceName": "wisdom_park"}, // 初始化字段，如：添加一个服务器名称
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, _ := config.Build()
	logger = logger.Named("")
	defer logger.Sync()
	logger.Info("log 初始化成功")
	logger.Info("无法获取网址",
		zap.String("url", "http://www.baidu.com"),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)
	//atom.SetLevel(zap.WarnLevel)
	//core := logging.Core().SetLevel
	core := logger.Core()
	fmt.Printf("%#v\n", core)
	logger.Core().Enabled(zapcore.Level(4))
	//= zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	//switch core.(type) {
	//case *zapcore.ioCore:
	//	()
	//}
	logger.Warn("无法获取网址",
		zap.String("url", "http://www.baidu.com"),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)

	logger.Error("无法获取网址",
		zap.String("url", "http://www.baidu.com"),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)

}

func TestMyLogger(t *testing.T) {
	logger := NewLogger("debug", "test")
	logger.Debug("debug")
	logger.Info("info")
	logger.Warn("warn")
	logger.Error("error")

	mylogger := NewSugaredLogger("debug", "printf")

	mylogger.Debugf("v =%s, u=%s", "key", "value")
	mylogger.Info("info")
	mylogger.Warn("warn")
	mylogger.Error("error")
}
