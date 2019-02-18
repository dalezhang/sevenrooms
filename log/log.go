package log

import (
	"go.uber.org/zap"
)

var (
	Logger  *zap.SugaredLogger
	ZLogger *zap.Logger
)

func Init(deubg bool) {
	var (
		err error
	)

	if deubg {
		ZLogger, err = zap.NewDevelopment()
	} else {
		ZLogger, err = zap.NewProduction()
	}
	if err != nil {
		panic(err)
	}

	Logger = ZLogger.Sugar()
}
