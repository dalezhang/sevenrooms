package logger

import (
	"os"
	"time"

	log "github.com/bigpigeon/logrus"
)

type NoneType struct{}

type FileLoggerInterval time.Duration

const (
	FileLoggerNone   = FileLoggerInterval(time.Duration(0))
	FileLoggerMinute = FileLoggerInterval(time.Minute)
	FileLoggerHour   = FileLoggerInterval(time.Hour)
	FileLoggerDay    = FileLoggerInterval(time.Hour * 24)
)

type FileLogger struct {
	log.Logger
	FilePath       string
	reopenSysSign  chan os.Signal
	reopenUserSign chan NoneType
	interval       FileLoggerInterval
	closeSign      chan NoneType
}

func CreateFileLogger(path string, interval FileLoggerInterval) (*FileLogger, error) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	logger := &FileLogger{
		Logger: log.Logger{
			Out:       f,
			Formatter: new(log.TextFormatter),
			Hooks:     make(log.LevelHooks),
			Level:     log.InfoLevel,
		},
		FilePath:       path,
		reopenSysSign:  make(chan os.Signal, 1),
		reopenUserSign: make(chan NoneType, 1),
		interval:       interval,
		closeSign:      make(chan NoneType, 1),
	}
	go logger.DoReopen()
	go logger.DoCron()

	return logger, nil
}

func (l *FileLogger) Reopen() error {
	l.Lock()
	defer l.Unlock()
	err := l.Out.(*os.File).Close()
	if err != nil {
		return err
	}
	f, err := os.OpenFile(l.FilePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	l.Out = f
	return nil
}

func (l *FileLogger) DoReopen() {

	for {
		select {
		case <-l.reopenSysSign:
			err := l.Reopen()
			if err != nil {
				panic(err)
			}
		case <-l.reopenUserSign:
			err := l.Reopen()
			if err != nil {
				panic(err)
			}
		case <-l.closeSign:
			return
		}
	}

}

func (l *FileLogger) ArchLogFile() {
	archFileName := l.FilePath + "." + time.Now().Format("20060102150405")
	os.Rename(l.FilePath, archFileName)
	l.reopenUserSign <- NoneType{}
}

func (l *FileLogger) DoCron() {
	if l.interval == FileLoggerNone {
		return
	}
	now := time.Now()
	firstWait := l.interval - FileLoggerInterval(now.UnixNano())%l.interval
	// square the time interval
	select {
	case <-time.After(time.Duration(firstWait)):
		l.ArchLogFile()
	case <-l.closeSign:
		return
	}
	for {
		select {
		case <-time.NewTicker(time.Duration(l.interval)).C:
			l.ArchLogFile()
		case <-l.closeSign:
			return
		}
	}
}

func (l *FileLogger) Close() {
	close(l.closeSign)
}
