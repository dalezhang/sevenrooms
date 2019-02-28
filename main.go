package main

import (
	"bindolabs/sevenrooms/config"
	"bindolabs/sevenrooms/db"
	"bindolabs/sevenrooms/log"
	"bindolabs/sevenrooms/models"
	"bindolabs/sevenrooms/module"
	"bindolabs/sevenrooms/sevenroom"
	"flag"
	"fmt"
	"time"
)

var (
	generate               = flag.Bool("g", false, "generate database tables")
	EventPartyExchangeName = "dale.local"
)

func main() {
	flag.Parse()
	err := config.Init()
	if err != nil {
		panic(err)
	}

	defer log.Logger.Sync()
	log.Logger.Infof("config init")
	if err = db.Init(log.Logger); err != nil {
		log.Logger.Errorf("init db err: %s", err)
		return
	}
	defer db.Exit()
	if *generate {
		log.Logger.Info("generate databases")
		if err := db.DB.LogMode(true).AutoMigrate(models.DBModels...).AutoMigrate(models.DBModels...).Error; err != nil {
			log.Logger.Error("generate databases err", err)
		}
		return
	}
	if err = sevenroom.Init(); err != nil {
		log.Logger.Errorf("init sevenrooms err: %s", err)
		return
	}
	GetPartyMsgFromDB()
}
func GetPartyMsgFromDB() {
	var lock bool
	go func() {
		i := 0
		for {
			if lock == false {
				lock = true
				fmt.Printf("\n ==========ScaneParties====================loop %d times\n", i)
				module.VoidChecks()
				err, rParties := module.ScaneParties()
				if err == nil {
					module.FilterPartyBySettingStores(rParties)
				}
				module.CloseChecks()
				module.CreateChecks()
				module.MoveTables()
				module.UpdateChecks()
				module.ListingVoidPayments()
				module.VoidPaymentOrCloseChecks()
				fmt.Printf("\n ==========sleep====================loop %d times\n", i)
				time.Sleep(time.Duration(10) * time.Second)
				lock = false
				i++
			}
		}
	}()

	forever := make(chan bool)
	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Logger.Errorf("%s: %s", msg, err)
	}
}
