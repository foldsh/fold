package logging

import (
	"errors"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
}

type LogLevel int

const (
	Fatal LogLevel = iota
	Error
	Warn
	Info
	Debug
)

func NewLogger(level LogLevel, json bool) (Logger, error) {
	var (
		config zap.Config
		logger *zap.Logger
		err    error
	)
	if json == true {
		config = zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(zapLevel(level))
		logger, err = config.Build()
	} else {
		config = zap.NewDevelopmentConfig()
		config.Level = zap.NewAtomicLevelAt(zapLevel(level))
		logger, err = config.Build()
	}
	if err != nil {
		return nil, errors.New("failed to create logger")
	}
	return logger.Sugar(), nil
}

func zapLevel(level LogLevel) zapcore.Level {
	switch level {
	case Fatal:
		return zapcore.FatalLevel
	case Error:
		return zapcore.ErrorLevel
	case Warn:
		return zapcore.WarnLevel
	case Info:
		return zapcore.InfoLevel
	case Debug:
		return zapcore.DebugLevel
	}
	return zapcore.InfoLevel
}

func NewTestLogger() Logger {
	return zap.NewExample().Sugar()
}
