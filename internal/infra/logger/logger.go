package logger

import (
	"os"
	"sync"
	"time"

	"github.com/felipeversiane/auth-service/internal/infra/config"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger      *zap.Logger
	once        sync.Once
	flushTicker *time.Ticker
	stopFlushCh chan struct{}
)

func New(cfg config.LogConfig) *zap.Logger {
	once.Do(func() {
		logLevel := getLogLevel(cfg.Level)

		var cores []zapcore.Core
		cores = append(cores, newFileCore(cfg, logLevel))

		if cfg.Environment == "development" {
			cores = append(cores, newConsoleCore(logLevel))
		}

		core := zapcore.NewTee(cores...)
		logger = zap.New(core).With(zap.String("service", cfg.ServiceName))

		startFlushRoutine()
	})

	return logger
}

func Info(msg string, fields ...zap.Field)  { log(zapcore.InfoLevel, msg, fields...) }
func Debug(msg string, fields ...zap.Field) { log(zapcore.DebugLevel, msg, fields...) }
func Warn(msg string, fields ...zap.Field)  { log(zapcore.WarnLevel, msg, fields...) }
func Error(msg string, fields ...zap.Field) { log(zapcore.ErrorLevel, msg, fields...) }
func Fatal(msg string, fields ...zap.Field) { log(zapcore.FatalLevel, msg, fields...) }

func Sync() {
	if logger != nil {
		_ = logger.Sync()
	}
}

func StopFlush() {
	if flushTicker != nil {
		flushTicker.Stop()
	}
	if stopFlushCh != nil {
		close(stopFlushCh)
	}
	Sync()
}

func log(level zapcore.Level, msg string, fields ...zap.Field) {
	if logger == nil {
		return
	}

	switch level {
	case zapcore.DebugLevel:
		logger.Debug(msg, fields...)
	case zapcore.InfoLevel:
		logger.Info(msg, fields...)
	case zapcore.WarnLevel:
		logger.Warn(msg, fields...)
	case zapcore.ErrorLevel:
		logger.Error(msg, fields...)
	case zapcore.FatalLevel:
		logger.Fatal(msg, fields...)
	}
}

func newFileCore(cfg config.LogConfig, logLevel zapcore.Level) zapcore.Core {
	logWriter := &lumberjack.Logger{
		Filename:   cfg.Path,
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     30,
		Compress:   true,
	}
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	return zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(logWriter),
		logLevel,
	)
}

func newConsoleCore(logLevel zapcore.Level) zapcore.Core {
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	return zapcore.NewCore(
		consoleEncoder,
		zapcore.AddSync(os.Stdout),
		logLevel,
	)
}

func startFlushRoutine() {
	flushTicker = time.NewTicker(5 * time.Second)
	stopFlushCh = make(chan struct{})

	go func() {
		for {
			select {
			case <-flushTicker.C:
				Sync()
			case <-stopFlushCh:
				return
			}
		}
	}()
}

func getLogLevel(level string) zapcore.Level {
	switch level {
	case "DEBUG":
		return zapcore.DebugLevel
	case "INFO":
		return zapcore.InfoLevel
	case "WARN":
		return zapcore.WarnLevel
	case "ERROR":
		return zapcore.ErrorLevel
	case "FATAL":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}
