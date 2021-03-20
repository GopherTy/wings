package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/gopherty/wings/common/conf"
	"github.com/gopherty/wings/utils"
)

var instance *zap.Logger

// Register 注册器
type Register struct {
}

// Name .
func (Register) Name() string {
	return "Common.Logger"
}

// Regist 实现 IRegister 接口，以注册获取初始化好的 logger 对象。
func (Register) Regist() (err error) {
	cnf := conf.Instance()

	// 是否输出日志文件
	var logPath []string
	if cnf.Logger.LogsPath != "" {
		// 创建指定路径
		err = utils.CreatePath(cnf.Logger.LogsPath)
		if err != nil {
			return
		}

		_, err = os.Create(cnf.Logger.LogsPath)
		if err != nil {
			return
		}

		logPath = []string{cnf.Logger.LogsPath}
	} else {
		logPath = []string{"stdout"}
	}

	var encodeTime zapcore.TimeEncoder
	if cnf.Logger.TimeFormat == "" {
		encodeTime = zapcore.ISO8601TimeEncoder
	} else {
		encodeTime = zapcore.TimeEncoderOfLayout(cnf.Logger.TimeFormat)
	}
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime:     encodeTime,                    // 默认为 ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 短径编码器
	}

	// 用户日志等级 debug,info,warn,error,dpanic,panic,fatal
	level := strings.TrimSpace(cnf.Logger.Level)
	level = strings.ToLower(level)
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	case "dpanic":
		zapLevel = zapcore.DPanicLevel
	case "panic":
		zapLevel = zapcore.PanicLevel
	case "fatal":
		zapLevel = zapcore.FatalLevel
	default:
		zapLevel = zapcore.DebugLevel
	}

	atomicLevel := zap.NewAtomicLevelAt(zapLevel)
	zapcnf := zap.Config{
		Level:            atomicLevel,
		Development:      cnf.Logger.Development,
		Encoding:         cnf.Logger.Encoding, // json 或 console
		OutputPaths:      logPath,
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig:    encoderConfig,
	}

	// 创建自定义日志对象
	zapLogger, err := zapcnf.Build()
	if err != nil {
		return
	}

	instance = zapLogger
	return
}

// Instance 获取默认的 logger 对象
func Instance() *zap.Logger {
	return instance
}
