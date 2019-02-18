package main

import (
	"bindolabs/optitable_middleware/config"
	"bindolabs/optitable_middleware/db"
	"bindolabs/optitable_middleware/log"
	"bindolabs/optitable_middleware/models"
	"bindolabs/optitable_middleware/module"
	"bindolabs/optitable_middleware/optitable"
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"github.com/streadway/amqp"
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
		log.Logger.Info("generate database tables")
		db.DB.LogMode(true).AutoMigrate(models.DBModels...).AutoMigrate(models.DBModels...)
		return
	}
	if err = optitable.Init(); err != nil {
		log.Logger.Errorf("init optitable err: %s", err)
		return
	}
	GetPartyMsgFromDB()
	// GetPartyMsgFromMQ()

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

func GetPartyMsgFromMQ() {
	fmt.Println("config.Conf.Setting.MQUrl = ", config.Conf.Setting.MQUrl)
	conn, err := amqp.Dial(config.Conf.Setting.MQUrl)

	defer conn.Close()
	if err != nil {
		log.Logger.Errorf("Failed to connect to RabbitMQ: %s", err)
		panic(err)
	}

	rabbitCh, err := conn.Channel()
	if err != nil {
		failOnError(err, "Failed to open a channel")
		panic(err)
	}

	defer rabbitCh.Close()

	q, err := rabbitCh.QueueDeclare(
		EventPartyExchangeName, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")
	err = rabbitCh.QueueBind(
		EventPartyExchangeName,
		"u.*",
		"parties",
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")
	if err != nil {
		panic(err)
	}

	msgs, err := rabbitCh.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			// str, _ := json.Marshal(d)
			// log.Logger.Infof("msg:", string(str))
			// log.Logger.Infof("Received a message: %s", d.Body)
			fmt.Printf("Received a message: %s\n", d.Body)
			var partyMessage module.PartyMqMessage
			json.Unmarshal(d.Body, &partyMessage)
			log.Logger.Infof("receive msg from Party: %s", partyMessage.PartyID)
			if err := partyMessage.FilterPartyBySettingStores(); err == nil {
				partyMessage.AnalyzeRecordFromMQ()
			}
			// hasGetMgs = true
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
