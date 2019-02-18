package env

import (
	"os"
	"path"

	"github.com/spf13/viper"
)

/**
Configure is a config manager work for low level(etc common lib,frameworks,...)
*/
type Configure struct {
	RunMode                string
	Version                string
	Services               map[string]Service
	Log                    Log
	SystemDatabase         *Database
	RubyBindoDatabase      *Database
	RubyGatewayDatabase    *Database
	RubyRestaurantDatabase *Database
	SystemCache            *Redis
	URLMap                 *URLMap
	DefaultSecret          *DefaultSecrets
	Secret                 *Secret
	Queue                  *QueueConfig
	Storage                *Storage
	Faye                   *Faye
	ElasticSearch          *ElasticSearch
	ETCD                   *Etcd
	Liquid                 *Liquid
	Hosts                  Hosts
	RabbitMQ               *RabbitMQ
}

func (c *Configure) ModuleENV(fileName string, v interface{}) error {
	envName := BindoCommonConfigFile
	if IsTest() {
		envName = BindoCommonTestConfigFile
	}
	vip := viper.New()
	vip.SetConfigType("yaml")
	vip.SetConfigFile(path.Dir(os.Getenv(envName)) + "/" + fileName)
	err := vip.ReadInConfig()
	if err != nil {
		panic(err)
	}
	return vip.Unmarshal(v)
}
