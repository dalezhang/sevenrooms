package logger

import (
	"os"

	log "github.com/bigpigeon/logrus"
)

var LogManager = map[string]*log.Logger{}

func Debug(v ...interface{}) {
	LogManager["default"].Debug(v...)
}

func Debugf(format string, v ...interface{}) {
	LogManager["default"].Debugf(format, v...)
}

func Info(v ...interface{}) {
	LogManager["default"].Info(v...)
}

func Infof(format string, v ...interface{}) {
	LogManager["default"].Infof(format, v...)
}

func Warn(v ...interface{}) {
	LogManager["default"].Warn(v...)
}

func Warnf(format string, v ...interface{}) {
	LogManager["default"].Warnf(format, v...)
}

func Error(v ...interface{}) {
	LogManager["default"].Error(v...)
}

func Errorf(format string, v ...interface{}) {
	LogManager["default"].Errorf(format, v...)
}

func init() {
	LogManager["default"] = &log.Logger{
		Out:       os.Stderr,
		Formatter: new(log.TextFormatter),
		Hooks:     make(log.LevelHooks),
		Level:     log.DebugLevel,
	}
}
