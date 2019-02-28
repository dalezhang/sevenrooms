package db

import (
	"bindolabs/sevenrooms/config"
	"flag"
	"fmt"
	"strings"

	"bindolabs/sevenrooms/log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go.uber.org/zap"
)

var (
	GatewayDB *gorm.DB
	RestDB    *gorm.DB
	DB        *gorm.DB
	showsql   = flag.Bool("showsql", false, "show sql info")
)

type Logger struct {
	*zap.SugaredLogger
}

func (l *Logger) Println(values ...interface{}) {
	flag.Parse()
	if *showsql {
		l.SugaredLogger.Info(values...)
	}
}
func Init(sl *zap.SugaredLogger) error {
	var err error
	for dbname, dbConfig := range config.Conf.DBs {
		fmt.Println("\n dbConfig.Database=====", dbConfig.Database)
		if strings.HasPrefix(dbConfig.Adapter, "mysql") {
			log.Logger.Debugf("%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local", dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database)
			switch dbname {
			case "gateway":
				GatewayDB, err = gorm.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local", dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database))
			case "restaurant":
				RestDB, err = gorm.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local", dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database))
			case "sevenroom":
				DB, err = gorm.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local", dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database))
			}
		} else {
			err = fmt.Errorf("not supported database adapter: %s", dbConfig.Adapter)
		}
		if err != nil {
			return err
		}
	}

	GatewayDB.SetLogger(gorm.Logger{
		LogWriter: &Logger{
			SugaredLogger: sl,
		},
	})
	RestDB.SetLogger(gorm.Logger{
		LogWriter: &Logger{
			SugaredLogger: sl,
		},
	})
	DB.SetLogger(gorm.Logger{
		LogWriter: &Logger{
			SugaredLogger: sl,
		},
	})

	if config.Conf.Debug {
		GatewayDB.LogMode(true)
		RestDB.LogMode(true)
		DB.LogMode(true)
	}

	return nil
}

func Exit() {
	if GatewayDB != nil {
		GatewayDB.Close()
	}
	if RestDB != nil {
		RestDB.Close()
	}
	if DB != nil {
		DB.Close()
	}
}
