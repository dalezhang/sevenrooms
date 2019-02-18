package config

import (
	"os"

	// "github.com/astaxie/beego/config"
	"github.com/spf13/viper"
)

var Config *viper.Viper

// func init() {
// 	Config = viper.New()
// 	Config.SetConfigType("yaml")
// 	Config.SetConfigFile(os.Getenv("BINDO_GATEWAY_CONFIG_FILE"))
// 	err := Config.ReadInConfig()
// 	if err != nil {
// 		panic(err)
// 	}
// }

func InitWithEnv(e string) {
	Config = viper.New()
	Config.SetConfigType("yaml")
	Config.SetConfigFile(os.Getenv(e))
	err := Config.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func NewConfigureFromEnv(env string, configType string) *viper.Viper {
	config := viper.New()
	config.SetConfigType(configType)
	config.SetConfigFile(os.Getenv(env))
	err := config.ReadInConfig()
	if err != nil {
		panic(err)
	}
	return config
}
