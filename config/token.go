package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func (c *Config) GetToken() string {
	c.lk.RLock()
	token := c.token
	c.lk.RUnlock()
	return token
}

func (c *Config) SetToken(token string) error {
	if token == "" {
		return fmt.Errorf("token is nil")
	}
	c.lk.Lock()
	c.token = token
	err := c.flushToken()
	c.lk.Unlock()
	return err
}

func (c *Config) flushToken() error {
	dir := filepath.Dir(c.tokenFile())
	if len(dir) > 0 && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return ioutil.WriteFile(c.tokenFile(), []byte(c.token), 0755)
}

func (c *Config) tokenFile() string {
	return filepath.Join("config", "token")
}
func (c *Config) loadToken() error {
	data, err := ioutil.ReadFile(c.tokenFile())
	if err != nil {
		return err
	}

	c.token = string(data)
	return nil
}
