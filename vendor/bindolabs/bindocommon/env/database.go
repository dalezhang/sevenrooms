package env

import "time"

type Database struct {
	DNS             string
	Driver          string
	ReadTimeout     int
	WriteTimeout    int
	MaxOpenConns    int
	ConnMaxLifeTime time.Duration
	MaxIdleConns    int
}
