package config

import (
	"bindolabs/sevenrooms/log"
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

	token string
	lk    sync.RWMutex

	DBs map[string]DB
	*Setting
}
type Setting struct {
	OpUrl  string           `yaml:"op_url"`
	Retry  int              `yaml:"retry"`
	Stores map[string]Store `yaml:"stores"`
	PosID  string           `yaml:"pos_id"`
	lk     sync.RWMutex
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
	StoreID      int    `yaml:"store_id"`
	VenueID      string `yaml:"venue_id"`
	Name         string `yaml:"name"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	Token        string
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
		Conf.Debug = true
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
	Conf.loadToken()
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

	if Conf.Setting.OpUrl == "" {
		errStrs = append(errStrs, "Config OpUrl is not set")
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
