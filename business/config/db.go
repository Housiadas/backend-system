package config

type DB struct {
	Name               string
	User               string
	Password           string
	Host               string
	MaxOpenConns       int
	MaxIdleConns       int
	ConnectionIdleTime string
	DisableTLS         bool
}
