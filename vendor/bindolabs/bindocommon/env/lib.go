package env

import (
	"fmt"
	"log"
	"os"
	"strings"

	"bindolabs/golib/config"
)

const (
	BindoCommonConfigFile     string = "BINDO_COMMON_CONFIG_FILE"      //the environment variable for the bindo common config
	BindoCommonTestConfigFile string = "BINDO_COMMON_TEST_CONFIG_FILE" //the environment variable for the bindo common config
	BindoCommonConfigType     string = "yaml"                          //we define the config use yaml format
	RunModeDev                string = "dev"                           //in your local machine or a dev cluster
	RunModeStaging            string = "staging"                       //in staging server
	RunModeProduction         string = "production"                    //in production server
)

var Env Configure

//is always auto call if environment variable `BINDO_COMMON_CONFIG_FILE` is set.
//else,will panic.
//nothing to do if environment variable `BINDO_COMMON_CONFIG_FILE`  not exists and its run in `go test`
func Init() {
	envName := BindoCommonConfigFile
	if IsTest() {
		envName = BindoCommonTestConfigFile
	}
	bindoCommonConfig := config.NewConfigureFromEnv(envName, BindoCommonConfigType)
	err := bindoCommonConfig.Unmarshal(&Env)
	if IsTest() {
		log.Printf("CommonConfigure:%+v \n", Env)
	}
	if err != nil {
		panic("Init Runtime Error:" + err.Error())
	}
}

//test it if run in staging environment
func IsStaging() bool {
	if Env.RunMode == RunModeStaging {
		return true
	}
	return false
}

//test it if run in production environment
func IsProduction() bool {
	if Env.RunMode == RunModeProduction {
		return true
	}
	return false
}

//test it if run in dev environment
func IsDev() bool {
	if Env.RunMode == RunModeDev {
		return true
	}
	return false
}

//return true if running in `go test`
//solution from: https://stackoverflow.com/a/45913089
func IsTest() bool {
	if strings.HasSuffix(os.Args[0], ".test") {
		return true
	}
	return false
}

func init() {
	if _, ok := os.LookupEnv(BindoCommonConfigFile); ok {
		Init()
	} else {
		if !IsTest() {
			panic(fmt.Sprintf("Evniroment Variable Not Found:%s", BindoCommonConfigFile))
		}
		log.Println("Evniroment Variable Not Found")
	}
}
