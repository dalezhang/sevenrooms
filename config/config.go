package config

import (
	"bindolabs/optitable_middleware/log"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	yaml "gopkg.in/yaml.v2"
)

const (
	Stg  = "staging"
	Prod = "production"
)

var Conf = new(Config)

type Config struct {
	HasInit bool
	Prod    bool
	Debug   bool

	lk sync.RWMutex

	DBs map[string]DB
	*Setting
}
type Setting struct {
	OpApiKey string           `yaml:"op_api_key"`
	OpUrl    string           `yaml:"op_url"`
	MQUrl    string           `yaml:"mq_url"`
	Retry    int              `yaml:"retry"`
	Stores   map[string]Store `yaml:"stores"`
}
type DB struct {
	Adapter  string
	Encoding string
	Database string
	Username string
	Password string
	Host     string
	Port     string
}
type Store struct {
	StoreID  int    `yaml:"store_id"`
	OpApiKey string `yaml:"op_api_key"`
}

func Init() error {
	if Conf.HasInit {
		return nil
	}
	// RUN_MODE=staging go run main.go
	runMode := os.Getenv("RUN_MODE")

	if runMode == Prod {
		fmt.Println("\n ENV=====production")
		log.Init(true)
		Conf.Prod = true
	} else {
		fmt.Println("\n ENV=====staging")
		log.Init(true)
		Conf.Debug = true
	}

	err := load(&Conf.DBs, "database.yml", runMode)
	if err != nil {
		return err
	}
	log.Logger.Debugf("DBs: %+v", Conf.DBs)

	Conf.Setting = new(Setting)
	if err = load(Conf.Setting, "setting.yml", runMode); err != nil {
		return err
	}
	ValidateSetting()
	Conf.HasInit = true
	return nil
}

func load(out interface{}, file, env string) error {
	b, err := ioutil.ReadFile(filepath.Join("config", file))
	if err != nil {
		return err
	}

	cs := make(map[string]interface{}, 4)
	err = yaml.Unmarshal(b, cs)
	if err != nil {
		return err
	}

	in, ok := cs[env]
	if !ok {
		return fmt.Errorf("cann't found env[%s] field in[%s]", env, file)
	}

	data, err := yaml.Marshal(in)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, out)
}

func ValidateSetting() {
	var errStrs []string
	if Conf.Setting.OpApiKey == "" {
		errStrs = append(errStrs, "Config OPApiKey is not set")
	}
	if Conf.Setting.OpUrl == "" {
		errStrs = append(errStrs, "Config OpUrl is not set")
	}
	if Conf.Setting.MQUrl == "" {
		errStrs = append(errStrs, "Config MQUrl is not set")
	}

	if len(errStrs) > 0 {
		fmt.Println(strings.Join(errStrs, "\n !!! "))
		panic(strings.Join(errStrs, ","))
	}
}

func GetStore(storeID int) (err error, store Store) {
	for _, store := range Conf.Setting.Stores {
		if store.StoreID == storeID {
			return nil, store
		}
	}
	err = fmt.Errorf("can't find store setting by store_id: %d", storeID)
	return
}
