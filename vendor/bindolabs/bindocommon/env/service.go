package env

import "errors"

type Service struct {
	Host        string
	Port        int
	Protocol    string
	Name        string
	DialTimeout int
	Timeout     int
}

func ServiceGetConfig(name string) (*Service, error) {
	c, ok := Env.Services[name]
	if !ok {
		return nil, errors.New("Service Not Found:" + name)
	}
	cc := new(Service)
	*cc = c
	return cc, nil
}
