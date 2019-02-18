package logger

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	log "github.com/bigpigeon/logrus"

	"bindolabs/bindocommon/helpers/gls"
	"github.com/jinzhu/gorm"
)

var LogManager = map[string]*log.Logger{}
var defaultLog *log.Logger

func getFileName() string {
	fileName := ""
	_, file, line, ok := runtime.Caller(2)
	if ok {
		fileName = fmt.Sprintf("%s:%d",file,line)
	}
	return fileName
}

func Debug(v ...interface{}) {
	GetLogger().WithField("file",getFileName()).WithField("mode", "debug").Debug(v...)
}

func Debugf(format string, v ...interface{}) {
	GetLogger().WithField("file",getFileName()).WithField("mode", "debug").Debugf(format, v...)
}

func Info(v ...interface{}) {
	GetLogger().WithField("file",getFileName()).WithField("mode", "info").Info(v...)
}

func Infof(format string, v ...interface{}) {
	GetLogger().WithField("file",getFileName()).WithField("mode", "info").Infof(format, v...)
}

func Warn(v ...interface{}) {
	GetLogger().WithField("file",getFileName()).WithField("mode", "warn").Warn(v...)
}

func Warnf(format string, v ...interface{}) {
	GetLogger().WithField("file",getFileName()).WithField("mode", "warn").Warnf(format, v...)
}

func Error(v ...interface{}) {
	GetLogger().WithField("mode", "error").Error(v...)
}

func Errorf(format string, v ...interface{}) {
	GetLogger().WithField("mode", "error").Errorf(format, v...)
}

func NewErrorf(format string, v ...interface{}) error {
	logs := fmt.Sprintf(format, v...)
	Error(logs)
	return errors.New(logs)
}

func WithFields(fields log.Fields) *log.Entry {
	logger := log.WithFields(fields)
	logger.Logger = &log.Logger{
		Out:       os.Stdout,
		Formatter: &log.JSONFormatter{},
		Hooks:     make(log.LevelHooks),
		Level:     log.DebugLevel,
	}
	return logger
}

func GetLogger() *log.Entry {
	data, ok := gls.Get("_log")
	if ok {
		return data.(*log.Entry)
	}
	return log.NewEntry(defaultLog)
}

func GoroutineInit(l *log.Entry) {
	gls.Set("_log", l)
}

func GoroutineShutdown() {
	gls.Shutdown()
}

func init() {
	loggerLevel := log.DebugLevel
	defaultLog = &log.Logger{
		Out:       os.Stdout,
		Formatter: &log.JSONFormatter{},
		Hooks:     make(log.LevelHooks),
		Level:     loggerLevel,
	}
}

type GormLogger struct {
	logInstance *log.Entry
}

func NewGormLogger() *GormLogger {
	return &GormLogger{logInstance: GetLogger()}
}

func (g *GormLogger) Print(v ...interface{}) {
	g.logInstance.WithFields(log.Fields{"module": "gorm", "type": "sql"}).Print(gorm.LogFormatter(v...)...)
}
