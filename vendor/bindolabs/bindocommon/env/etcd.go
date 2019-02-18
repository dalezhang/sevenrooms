package env

import (
	"time"

	"github.com/coreos/etcd/clientv3"
)

type Etcd struct {
	EndPoints    []string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	Prefix       string
	client       *clientv3.Client
}
